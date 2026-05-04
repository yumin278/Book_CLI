package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"molt/internal/moltbook"
)

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}

	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 3 {
		return string(r[:max])
	}
	return string(r[:max-3]) + "..."
}

func formatStatus(res *moltbook.StatusResponse) error {
	fmt.Printf("🦞 Your Moltbook Profile\n")
	fmt.Printf("Name: %s\n", res.Agent.Name)
	fmt.Printf("Status: %s\n", res.Status)
	fmt.Printf("Karma: %d\n", res.Agent.Karma)
	fmt.Printf("Description: %s\n", res.Agent.Description)
	return nil
}

func printPostsTable(posts []moltbook.Post) error {
	if len(posts) == 0 {
		fmt.Println("No posts found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSubmolt\tTitle\tAuthor\tUpvotes\tComments")
	for _, p := range posts {
		id := p.ID
		title := truncateRunes(p.Title, 40)
		submolt := p.Submolt.Name
		if submolt == "" {
			submolt = "?"
		}
		author := p.Author.Name
		if author == "" {
			author = "?"
		}
		fmt.Fprintf(w, "%s\tm/%s\t%s\t%s\t%d\t%d\n", id, submolt, title, author, p.Upvotes, p.CommentCount)
	}
	return w.Flush()
}

func formatFeed(res *moltbook.FeedResponse) error {
	return printPostsTable(res.Posts)
}

func formatPostCreated(res *moltbook.PostResponse) error {
	fmt.Printf("✓ Post created! ID: %s\n", res.Post.ID)
	return nil
}

func formatCommentsTable(comments []moltbook.Comment, full bool) error {
	if len(comments) == 0 {
		fmt.Println("No comments found.")
		return nil
	}

	if full {
		for _, c := range comments {
			author := c.Author.Name
			if author == "" {
				author = "?"
			}
			fmt.Printf("Author: %s | Score: %d | Replies: %d\n", author, c.Score, c.ReplyCount)
			fmt.Printf("%s\n\n", c.Content)
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Author\tScore\tReplies\tContent")
	for _, c := range comments {
		author := c.Author.Name
		if author == "" {
			author = "?"
		}

		content := c.Content
		// Remove newlines for table formatting
		content = strings.ReplaceAll(content, "\n", " ")
		content = truncateRunes(content, 60)

		fmt.Fprintf(w, "%s\t%d\t%d\t%s\n", author, c.Score, c.ReplyCount, content)
	}
	return w.Flush()
}

func formatCommentsList(res *moltbook.CommentsResponse, full bool, postID string) error {
	fmt.Printf("Top %d comments:\n", len(res.Comments))
	err := formatCommentsTable(res.Comments, full)
	if err != nil {
		return err
	}
	if res.NextCursor != "" {
		fmt.Printf("\nNext: molt post comments %s --cursor %s\n", postID, res.NextCursor)
	}
	return nil
}

func formatDMCheck(res *moltbook.DMCheckResponse) error {
	if res.HasActivity {
		summary := res.Summary
		if summary == "" {
			summary = "You have activity!"
		}
		fmt.Printf("📬 %s\n", summary)
	} else {
		fmt.Println("No new DM activity.")
	}
	return nil
}

func formatDMConversations(res *moltbook.DMConversationsResponse) error {
	convos := res.Conversations.Items
	if len(convos) == 0 {
		fmt.Println("No conversations yet.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tWith\tUnread\tLast Activity")
	for _, c := range convos {
		id := c.ConversationID
		with := c.WithAgent.Name
		if with == "" {
			with = "?"
		}
		time := c.LastMessageAt
		if len(time) > 10 {
			time = time[:10]
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", id, with, c.UnreadCount, time)
	}
	return w.Flush()
}

func formatDMRead(res *moltbook.DMReadResponse) error {
	for _, msg := range res.Messages {
		sender := msg.From.Name
		if sender == "" {
			sender = "?"
		}
		time := msg.CreatedAt
		if len(time) > 16 {
			time = time[:16]
		}
		fmt.Printf("%s (%s)\n", sender, time)
		fmt.Printf("  %s\n\n", msg.Message)
	}
	return nil
}

func formatDMRequests(res *moltbook.DMRequestsResponse) error {
	if len(res.Requests) == 0 {
		fmt.Println("No pending requests.")
		return nil
	}

	for _, req := range res.Requests {
		from := req.From.Name
		if from == "" {
			from = "?"
		}
		preview := truncateRunes(req.MessagePreview, 53)
		convID := req.ConversationID
		fmt.Printf("%s (%s)\n", from, convID)
		fmt.Printf("  %s\n\n", preview)
	}
	return nil
}

func formatSubmolts(res *moltbook.SubmoltsResponse) error {
	if len(res.Submolts) == 0 {
		fmt.Println("No submolts found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tDescription\tMembers")
	for _, sub := range res.Submolts {
		name := sub.Name
		if name == "" {
			name = "?"
		}
		desc := truncateRunes(sub.Description, 40)
		fmt.Fprintf(w, "m/%s\t%s\t%d\n", name, desc, sub.MemberCount)
	}
	return w.Flush()
}

func formatSearch(res *moltbook.SearchResponse) error {
	if len(res.Results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	for _, item := range res.Results {
		itemType := item.Type
		if itemType == "" {
			itemType = "post"
		}

		title := item.Title
		if title == "" {
			title = item.Content
		}
		title = truncateRunes(title, 50)

		fmt.Printf("[%s] %s\n", itemType, title)
	}
	return nil
}

func formatSuccessMessage(msg string) func() error {
	return func() error {
		fmt.Println(msg)
		return nil
	}
}

func formatPost(res *moltbook.PostResponse) error {
	p := res.Post
	fmt.Printf("[%s] %s\n", p.ID, p.Title)
	if p.URL != "" {
		fmt.Printf("URL: %s\n", p.URL)
	}
	if p.Content != "" {
		fmt.Printf("\n%s\n", p.Content)
	}
	fmt.Printf("\nAuthor: %s | Submolt: m/%s | ⬆ %d | 💬 %d\n", p.Author.Name, p.Submolt.Name, p.Upvotes, p.CommentCount)
	return nil
}

func formatCommentCreated(res *moltbook.CommentCreateResponse) error {
	fmt.Printf("✓ Comment added! ID: %s\n", res.Comment.ID)
	return nil
}

func formatHome(res *moltbook.HomeResponse) error {
	fmt.Printf("🏠 Welcome back, %s! (Karma: %d)\n", res.YourAccount.Name, res.YourAccount.Karma)
	fmt.Printf("Notifications: %d unread\n", res.YourAccount.UnreadNotificationCount)
	fmt.Printf("DMs: %s unread, %s pending requests\n\n", res.YourDirectMessages.UnreadMessageCount, res.YourDirectMessages.PendingRequestCount)

	if len(res.ActivityOnYourPosts) > 0 {
		fmt.Printf("🔔 Activity on your posts:\n")
		for _, act := range res.ActivityOnYourPosts {
			fmt.Printf("- Post [%s] (%d new notifications)\n", act.PostID, act.NewNotificationCount)
			for _, action := range act.SuggestedActions {
				fmt.Printf("  > %s\n", translateSuggestedAction(action))
			}
		}
		fmt.Println()
	}

	if len(res.PostsFromAccountsYouFollow.Posts) > 0 {
		fmt.Printf("📝 From accounts you follow (Total: %d):\n", res.PostsFromAccountsYouFollow.TotalFollowing)
		for _, p := range res.PostsFromAccountsYouFollow.Posts {
			fmt.Printf("- [%s] m/%s | %s\n", p.PostID, p.SubmoltName, p.Title)
			fmt.Printf("  By: %s | ⬆ %d | 💬 %d  (Use: molt post get %s)\n", p.AuthorName, p.Upvotes, p.CommentCount, p.PostID)
		}
		fmt.Println()
	}

	return nil
}
