package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"

	"github.com/Sush1sui/thesis-bot-pinger-go/internal/bot"
	"github.com/bwmarrin/discordgo"
)

var devDiscordIDs = []string{
	"982491279369830460",
	"990103932669931580",
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Thesis Bot Pinger!"))
}

type Commit struct {
	URL string `json:"url"`
	Message string `json:"message"`
	Author struct {
		Username string `json:"username"`
	} `json:"author"`
}

type Repository struct {
	Name string `json:"name"`
	FullName string `json:"full_name"`
}

type Sender struct {
	AvatarURL string `json:"avatar_url"`
}

type Payload struct {
	Ref string `json:"ref"`
	HeadCommit Commit `json:"head_commit"`
	Sender Sender `json:"sender"`
	Repository Repository `json:"repository"`
}

func SendNotification(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Send Notification called")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Received GitHub webhook request")

	var payload Payload
	var body []byte
	var err error

	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		// Parse form and extract payload
		if err := r.ParseForm(); err != nil {
			fmt.Println("Error parsing form:", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		payloadStr := r.FormValue("payload")
		if payloadStr == "" {
			fmt.Println("No payload found in form")
			http.Error(w, "No payload found", http.StatusBadRequest)
			return
		}
		body = []byte(payloadStr)
	} else {
		body, err = io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil || len(body) == 0 {
			fmt.Println("Error reading request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Println("Received webhook payload:", payload)

	if payload.Ref == "refs/heads/main" || payload.Ref == "refs/heads/master" {
		commit := payload.HeadCommit
		ogImage, err := getOpenGraphImage(commit.URL)
		if err != nil {
			ogImage = payload.Sender.AvatarURL
			fmt.Println("Error fetching Open Graph image:", err)
		}

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("New Commit by %s on Repo: %s", commit.Author.Username, payload.Repository.Name),
			Color: 0xFFFFFF,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "Full Repository Name",
					Value: payload.Repository.FullName,
					Inline: true,
				},
				{
					Name: "Author",
					Value: commit.Author.Username,
					Inline: true,
				},
				{
					Name: "Commit Message",
					Value: commit.Message,
					Inline: false,
				},
				{
					Name: "Commit URL",
					Value: "[View Commit](" + commit.URL + ")",
					Inline: false,
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: ogImage,
			},
		}

		go func() {
			_, err := bot.Session.ChannelMessageSendComplex("1373128370358980770", &discordgo.MessageSend{
				Content: "**Hello <@&1350439469462978620>! There is a new commit!**",
				Embed: embed,
				AllowedMentions: &discordgo.MessageAllowedMentions{Roles: []string{"1350439469462978620"}},
			})
			if err != nil {
				fmt.Println("Error sending Discord message:", err)
			}
		}()

		for _, id := range devDiscordIDs {
			go func() {
				user, err := bot.Session.User(id)
				if err != nil || user == nil {
					fmt.Println("Error fetching user:", err)
					return
				}
				dmChannel, err := bot.Session.UserChannelCreate(user.ID)
				if err != nil {
					fmt.Println("Error creating DM channel:", err)
					return
				}

				_, err = bot.Session.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
					Content: fmt.Sprintf("**Hello <@%s>! There is a new commit!**", user.ID),
					Embed: embed,
					AllowedMentions: &discordgo.MessageAllowedMentions{Roles: []string{"1350439469462978620"}},
				})
				if err != nil {
					fmt.Println("Error sending DM message:", err)
				}
			}()
		}

		fmt.Println("Notification sent successfully")
	}
}

// func verifyGithubSignature(r *http.Request, body []byte) error {
// 	signature := r.Header.Get("X-Hub-Signature-256")
// 	if signature == "" {
// 		return fmt.Errorf("missing signature header")
// 	}

// 	secret := config.GlobalConfig.GithubSecret
// 	if secret == "" {
// 		return fmt.Errorf("missing GitHub secret in configuration")
// 	}

// 	mac := hmac.New(sha256.New, []byte(secret))
// 	mac.Write(body)
// 	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
// 	if !hmac.Equal([]byte(signature), []byte(expected)) {
// 		return fmt.Errorf("invalid signature: expected %s, got %s", expected, signature)
// 	}
// 	return nil
// }

func getOpenGraphImage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var ogImage string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var property, content string
			for _, attr := range n.Attr {
				if attr.Key == "property" && attr.Val == "og:image" {
					property = attr.Val
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}
			if property == "og:image" && content != "" {
				ogImage = content
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if ogImage == "" {
		return "", fmt.Errorf("no Open Graph image found")
	}
	return ogImage, nil
}