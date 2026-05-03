package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"molt/internal/moltbook"
)

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
		if len(id) > 8 {
			id = id[:8]
		}
		title := p.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}
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
		if len(id) > 8 {
			id = id[:8]
		}
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
		preview := req.MessagePreview
		if len(preview) > 50 {
			preview = preview[:50] + "..."
		}
		convID := req.ConversationID
		if len(convID) > 8 {
			convID = convID[:8]
		}
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
		desc := sub.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
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
		if len(title) > 50 {
			title = title[:47] + "..."
		}

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

func formatCommentCreated(res *moltbook.SuccessResponse) error {
	fmt.Println("✓ Comment posted!")
	return nil
}
