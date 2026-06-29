package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"url-shortner/internal/storage"
)

func RunInteractiveCLI(store storage.Storage) {

	helpStr := `
Commands:
  help                 Show help
  shorten <url>        Create short URL
  get <code>           Get original URL
  delete <code>        Delete URL
  list                 List URLs
  exit                 Quit
`

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to URL Shortner")
	fmt.Println(helpStr)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		args := strings.Fields(strings.TrimSpace(line))

		if len(args) == 0 {
			continue
		}

		cmd := args[0]

		switch cmd {

		case "help":
			fmt.Println(helpStr)

		case "shorten":
			if len(args) != 2 {
				fmt.Println("usage: shorten <url>")
				continue
			}

			url, err := store.CreateURL(args[1])
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			fmt.Println("created:", url)

		case "get":
			if len(args) != 2 {
				fmt.Println("usage: get <code>")
				continue
			}

			url, err := store.GetOriginalURLById(args[1])
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			fmt.Println("redirect:", url)

		case "delete":
			if len(args) != 2 {
				fmt.Println("usage: delete <code>")
				continue
			}

			err := store.DeleteURL(args[1])
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			fmt.Println("deleted:", args[1])

		case "list":
			urls, err := store.GetURLs()
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

			fmt.Fprintln(w, "CODE\tURL\tCREATED")
			fmt.Fprintln(w, "----\t---\t-------")

			for _, item := range urls {
				fmt.Fprintf(
					w,
					"%s\t%s\t%s\n",
					item.Id,
					item.RedirectTO,
					item.CreatedAt,
				)
			}

			w.Flush()

		case "exit", "quit":
			fmt.Println("Bye!")
			return

		default:
			fmt.Println("unknown command. type 'help'")
		}
	}
}
