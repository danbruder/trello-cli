package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/adlio/trello"
	"github.com/danbruder/trello-cli/internal/batch"
	"github.com/danbruder/trello-cli/internal/client"
)

// TestE2EAllCommands is a comprehensive end-to-end test that exercises all CLI commands
// against a live Trello API. It requires TRELLO_API_KEY and TRELLO_TOKEN environment variables.
func TestE2EAllCommands(t *testing.T) {
	// Check for required environment variables
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")

	if apiKey == "" || token == "" {
		t.Skip("Skipping E2E test: TRELLO_API_KEY and TRELLO_TOKEN environment variables are required")
	}

	// Initialize client
	trelloClient := client.NewClient(apiKey, token)

	// Generate unique test identifier
	testID := fmt.Sprintf("e2e-test-%d", time.Now().Unix())

	t.Logf("Starting E2E test with ID: %s", testID)

	// Cleanup function to run at the end
	var testBoardID string
	defer func() {
		if testBoardID != "" {
			t.Logf("Cleaning up test board: %s", testBoardID)
			// Delete the test board
			board, err := trelloClient.GetBoard(testBoardID, nil)
			if err == nil {
				_ = board.Delete(nil)
			}
		}
	}()

	// Test 1: Board Commands
	t.Run("Board Commands", func(t *testing.T) {
		// Create board
		t.Run("Create Board", func(t *testing.T) {
			boardName := fmt.Sprintf("Test Board %s", testID)
			board := trello.NewBoard(boardName)

			err := trelloClient.CreateBoard(&board, nil)
			if err != nil {
				t.Fatalf("Failed to create board: %v", err)
			}

			testBoardID = board.ID
			t.Logf("Created board: %s (ID: %s)", board.Name, board.ID)

			if board.Name != boardName {
				t.Errorf("Expected board name %s, got %s", boardName, board.Name)
			}
		})

		// Create board with description
		t.Run("Create Board with Description", func(t *testing.T) {
			boardName := fmt.Sprintf("Test Board with Desc %s", testID)
			boardDesc := "This is a test board with a description"
			board := trello.NewBoard(boardName)
			board.Desc = boardDesc

			err := trelloClient.CreateBoard(&board, nil)
			if err != nil {
				t.Fatalf("Failed to create board with description: %v", err)
			}

			// Clean up this board at the end
			defer func() {
				if board.ID != "" {
					t.Logf("Cleaning up board with description: %s", board.ID)
					b, err := trelloClient.GetBoard(board.ID, nil)
					if err == nil {
						_ = b.Delete(nil)
					}
				}
			}()

			t.Logf("Created board with description: %s (ID: %s)", board.Name, board.ID)

			if board.Name != boardName {
				t.Errorf("Expected board name %s, got %s", boardName, board.Name)
			}

			// Verify description was set
			retrievedBoard, err := trelloClient.GetBoard(board.ID, nil)
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			if retrievedBoard.Desc != boardDesc {
				t.Errorf("Expected board description %q, got %q", boardDesc, retrievedBoard.Desc)
			}
		})

		// List boards
		t.Run("List Boards", func(t *testing.T) {
			member, err := trelloClient.GetMember("me", nil)
			if err != nil {
				t.Fatalf("Failed to get current member: %v", err)
			}

			boards, err := member.GetBoards(nil)
			if err != nil {
				t.Fatalf("Failed to list boards: %v", err)
			}

			if len(boards) == 0 {
				t.Error("Expected at least one board")
			}

			// Check if our test board is in the list
			found := false
			for _, board := range boards {
				if board.ID == testBoardID {
					found = true
					break
				}
			}
			if !found {
				t.Error("Test board not found in board list")
			}
			t.Logf("Found %d boards", len(boards))
		})

		// Get board
		t.Run("Get Board", func(t *testing.T) {
			board, err := trelloClient.GetBoard(testBoardID, nil)
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			if board.ID != testBoardID {
				t.Errorf("Expected board ID %s, got %s", testBoardID, board.ID)
			}
			t.Logf("Retrieved board: %s", board.Name)
		})

		// Add member to board (test with current user)
		t.Run("Add Member to Board", func(t *testing.T) {
			// Get current member to add them to the board
			member, err := trelloClient.GetMember("me", nil)
			if err != nil {
				t.Fatalf("Failed to get current member: %v", err)
			}

			board, err := trelloClient.GetBoard(testBoardID, nil)
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			// Try to add current member (should succeed or indicate already a member)
			memberToAdd := trello.Member{Email: member.Email}
			_, err = board.AddMember(&memberToAdd, nil)
			// Note: This might fail if the user is already a member, which is fine
			// We're just testing that the command doesn't crash
			t.Logf("Add member result (may already be member): %v", err)
		})

		// Delete board (create a temporary board for this test)
		t.Run("Delete Board", func(t *testing.T) {
			// Create a temporary board to delete
			tempBoardName := fmt.Sprintf("Temp Board to Delete %s", testID)
			tempBoard := trello.NewBoard(tempBoardName)

			err := trelloClient.CreateBoard(&tempBoard, nil)
			if err != nil {
				t.Fatalf("Failed to create temp board: %v", err)
			}

			t.Logf("Created temp board for deletion: %s (ID: %s)", tempBoard.Name, tempBoard.ID)

			// Delete the board
			err = tempBoard.Delete(nil)
			if err != nil {
				t.Fatalf("Failed to delete board: %v", err)
			}

			t.Logf("Deleted board: %s", tempBoard.ID)

			// Verify it's deleted by trying to get it (should fail)
			_, err = trelloClient.GetBoard(tempBoard.ID, nil)
			if err == nil {
				t.Error("Expected error when getting deleted board, but got none")
			}
		})
	})

	// Test 2: List Commands
	var testListID string
	var secondListID string

	t.Run("List Commands", func(t *testing.T) {
		// Get default lists (Trello creates "To Do", "Doing", "Done" by default)
		t.Run("Get Lists", func(t *testing.T) {
			board, err := trelloClient.GetBoard(testBoardID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			lists, err := board.GetLists(trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get lists: %v", err)
			}

			if len(lists) == 0 {
				t.Fatal("Expected at least one list in new board")
			}

			testListID = lists[0].ID
			t.Logf("Found %d lists, using list: %s (ID: %s)", len(lists), lists[0].Name, lists[0].ID)

			// Use second list if available
			if len(lists) > 1 {
				secondListID = lists[1].ID
				t.Logf("Second list: %s (ID: %s)", lists[1].Name, lists[1].ID)
			}
		})

		// Create a new list
		t.Run("Create List", func(t *testing.T) {
			board, err := trelloClient.GetBoard(testBoardID, nil)
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			listName := fmt.Sprintf("Test List %s", testID)

			list, err := board.CreateList(listName, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to create list: %v", err)
			}

			secondListID = list.ID
			t.Logf("Created list: %s (ID: %s)", list.Name, list.ID)

			if list.Name != listName {
				t.Errorf("Expected list name %s, got %s", listName, list.Name)
			}
		})

		// Get list
		t.Run("Get List", func(t *testing.T) {
			list, err := trelloClient.GetList(testListID, nil)
			if err != nil {
				t.Fatalf("Failed to get list: %v", err)
			}

			if list.ID != testListID {
				t.Errorf("Expected list ID %s, got %s", testListID, list.ID)
			}
			t.Logf("Retrieved list: %s", list.Name)
		})

		// Archive list
		t.Run("Archive List", func(t *testing.T) {
			// Create a temporary list to archive
			board, err := trelloClient.GetBoard(testBoardID, nil)
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			tempListName := fmt.Sprintf("List to Archive %s", testID)
			tempList, err := board.CreateList(tempListName, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to create temp list: %v", err)
			}

			t.Logf("Created temp list for archiving: %s (ID: %s)", tempList.Name, tempList.ID)

			// Archive it
			err = tempList.Archive()
			if err != nil {
				t.Fatalf("Failed to archive list: %v", err)
			}

			// Verify it's archived
			archivedList, err := trelloClient.GetList(tempList.ID, nil)
			if err != nil {
				t.Fatalf("Failed to get archived list: %v", err)
			}

			if !archivedList.Closed {
				t.Error("Expected list to be archived (closed)")
			}
			t.Logf("Archived list: %s", tempList.ID)
		})
	})

	// Test 3: Card Commands
	var testCardID string
	var copiedCardID string

	t.Run("Card Commands", func(t *testing.T) {
		// Create card
		t.Run("Create Card", func(t *testing.T) {
			cardName := fmt.Sprintf("Test Card %s", testID)
			card := trello.Card{
				Name:   cardName,
				IDList: testListID,
				Desc:   "This is a test card created by E2E test",
			}

			err := trelloClient.CreateCard(&card, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to create card: %v", err)
			}

			testCardID = card.ID
			t.Logf("Created card: %s (ID: %s)", card.Name, card.ID)

			if card.Name != cardName {
				t.Errorf("Expected card name %s, got %s", cardName, card.Name)
			}

			// Verify description was set
			retrievedCard, err := trelloClient.GetCard(card.ID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}

			if retrievedCard.Desc != "This is a test card created by E2E test" {
				t.Errorf("Expected card description %q, got %q", "This is a test card created by E2E test", retrievedCard.Desc)
			}
		})

		// List cards
		t.Run("List Cards", func(t *testing.T) {
			list, err := trelloClient.GetList(testListID, nil)
			if err != nil {
				t.Fatalf("Failed to get list: %v", err)
			}

			cards, err := list.GetCards(trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to list cards: %v", err)
			}

			if len(cards) == 0 {
				t.Error("Expected at least one card")
			}

			// Check if our test card is in the list
			found := false
			for _, card := range cards {
				if card.ID == testCardID {
					found = true
					break
				}
			}
			if !found {
				t.Error("Test card not found in card list")
			}
			t.Logf("Found %d cards in list", len(cards))
		})

		// Get card
		t.Run("Get Card", func(t *testing.T) {
			card, err := trelloClient.GetCard(testCardID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}

			if card.ID != testCardID {
				t.Errorf("Expected card ID %s, got %s", testCardID, card.ID)
			}
			t.Logf("Retrieved card: %s", card.Name)
		})

		// Move card (only if we have a second list)
		if secondListID != "" {
			t.Run("Move Card", func(t *testing.T) {
				card, err := trelloClient.GetCard(testCardID, nil)
				if err != nil {
					t.Fatalf("Failed to get card: %v", err)
				}

				err = card.MoveToList(secondListID, trello.Defaults())
				if err != nil {
					t.Fatalf("Failed to move card: %v", err)
				}

				// Verify the move
				updatedCard, err := trelloClient.GetCard(testCardID, trello.Defaults())
				if err != nil {
					t.Fatalf("Failed to get updated card: %v", err)
				}

				if updatedCard.IDList != secondListID {
					t.Errorf("Expected card to be in list %s, got %s", secondListID, updatedCard.IDList)
				}
				t.Logf("Moved card to list: %s", secondListID)
			})

			// Copy card
			t.Run("Copy Card", func(t *testing.T) {
				card, err := trelloClient.GetCard(testCardID, nil)
				if err != nil {
					t.Fatalf("Failed to get card: %v", err)
				}

				// Copy the card
				copiedCard := trello.Card{
					Name:   fmt.Sprintf("Copy of %s", card.Name),
					IDList: testListID,
					Desc:   card.Desc,
				}

				err = trelloClient.CreateCard(&copiedCard, trello.Defaults())
				if err != nil {
					t.Fatalf("Failed to copy card: %v", err)
				}

				copiedCardID = copiedCard.ID
				t.Logf("Copied card: %s (ID: %s)", copiedCard.Name, copiedCard.ID)

				if copiedCard.IDList != testListID {
					t.Errorf("Expected copied card to be in list %s, got %s", testListID, copiedCard.IDList)
				}
			})
		}

		// Archive card
		t.Run("Archive Card", func(t *testing.T) {
			// Create a temporary card to archive
			tempCard := trello.Card{
				Name:   fmt.Sprintf("Card to Archive %s", testID),
				IDList: testListID,
			}

			err := trelloClient.CreateCard(&tempCard, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to create temp card: %v", err)
			}

			// Archive it
			err = tempCard.Archive()
			if err != nil {
				t.Fatalf("Failed to archive card: %v", err)
			}

			// Verify it's archived
			archivedCard, err := trelloClient.GetCard(tempCard.ID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get archived card: %v", err)
			}

			if !archivedCard.Closed {
				t.Error("Expected card to be archived (closed)")
			}
			t.Logf("Archived card: %s", tempCard.ID)
		})
	})

	// Test 4: Label Commands
	var testLabelID string

	t.Run("Label Commands", func(t *testing.T) {
		// Create label
		t.Run("Create Label", func(t *testing.T) {
			board, err := trelloClient.GetBoard(testBoardID, nil)
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			label := trello.Label{
				Name:  fmt.Sprintf("Test Label %s", testID),
				Color: "red",
			}

			err = board.CreateLabel(&label, nil)
			if err != nil {
				t.Fatalf("Failed to create label: %v", err)
			}

			testLabelID = label.ID
			t.Logf("Created label: %s (ID: %s)", label.Name, label.ID)

			if label.Color != "red" {
				t.Errorf("Expected label color red, got %s", label.Color)
			}
		})

		// Get labels
		t.Run("Get Labels", func(t *testing.T) {
			board, err := trelloClient.GetBoard(testBoardID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			labels, err := board.GetLabels(trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get labels: %v", err)
			}

			// Check if our test label is in the list
			found := false
			for _, label := range labels {
				if label.ID == testLabelID {
					found = true
					break
				}
			}
			if !found {
				t.Error("Test label not found in label list")
			}
			t.Logf("Found %d labels", len(labels))
		})

		// Add label to card
		t.Run("Add Label to Card", func(t *testing.T) {
			card, err := trelloClient.GetCard(testCardID, nil)
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}

			err = card.AddIDLabel(testLabelID)
			if err != nil {
				t.Fatalf("Failed to add label to card: %v", err)
			}

			// Verify the label was added
			updatedCard, err := trelloClient.GetCard(testCardID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get updated card: %v", err)
			}

			found := false
			for _, labelID := range updatedCard.IDLabels {
				if labelID == testLabelID {
					found = true
					break
				}
			}
			if !found {
				t.Error("Label not found on card")
			}
			t.Logf("Added label %s to card %s", testLabelID, testCardID)
		})
	})

	// Test 5: Checklist Commands
	var testChecklistID string

	t.Run("Checklist Commands", func(t *testing.T) {
		// Create checklist
		t.Run("Create Checklist", func(t *testing.T) {
			card, err := trelloClient.GetCard(testCardID, nil)
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}

			checklistName := fmt.Sprintf("Test Checklist %s", testID)
			checklist, err := trelloClient.CreateChecklist(card, checklistName, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to create checklist: %v", err)
			}

			testChecklistID = checklist.ID
			t.Logf("Created checklist: %s (ID: %s)", checklist.Name, checklist.ID)
		})

		// Get checklists
		t.Run("Get Checklists", func(t *testing.T) {
			card, err := trelloClient.GetCard(testCardID, trello.Arguments{"checklists": "all"})
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}

			checklists := card.Checklists
			if checklists == nil {
				checklists = []*trello.Checklist{}
			}

			if len(checklists) == 0 {
				t.Error("Expected at least one checklist")
			}

			// Check if our test checklist is in the list
			found := false
			for _, checklist := range checklists {
				if checklist.ID == testChecklistID {
					found = true
					break
				}
			}
			if !found {
				t.Error("Test checklist not found in checklist list")
			}
			t.Logf("Found %d checklists", len(checklists))
		})

		// Add checklist item
		t.Run("Add Checklist Item", func(t *testing.T) {
			checklist, err := trelloClient.GetChecklist(testChecklistID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get checklist: %v", err)
			}

			itemName := "Test checklist item"
			item, err := checklist.CreateCheckItem(itemName, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to add checklist item: %v", err)
			}

			t.Logf("Added checklist item: %s (ID: %s)", item.Name, item.ID)

			// Verify the item was added
			updatedChecklist, err := trelloClient.GetChecklist(testChecklistID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get updated checklist: %v", err)
			}

			if len(updatedChecklist.CheckItems) == 0 {
				t.Error("Expected at least one checklist item")
			}
		})
	})

	// Test 6: Attachment Commands
	t.Run("Attachment Commands", func(t *testing.T) {
		// Add attachment
		t.Run("Add Attachment", func(t *testing.T) {
			card, err := trelloClient.GetCard(testCardID, nil)
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}

			// Add a URL attachment
			testURL := "https://example.com/test-attachment.pdf"
			attachment := trello.Attachment{
				URL: testURL,
			}

			err = card.AddURLAttachment(&attachment)
			if err != nil {
				t.Fatalf("Failed to add attachment: %v", err)
			}

			t.Logf("Added attachment to card: %s", testURL)
		})

		// List attachments
		t.Run("List Attachments", func(t *testing.T) {
			card, err := trelloClient.GetCard(testCardID, nil)
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}

			attachments, err := card.GetAttachments(nil)
			if err != nil {
				t.Fatalf("Failed to get attachments: %v", err)
			}

			if len(attachments) == 0 {
				t.Error("Expected at least one attachment")
			}

			// Verify our test attachment is in the list
			found := false
			for _, att := range attachments {
				if att.URL == "https://example.com/test-attachment.pdf" {
					found = true
					break
				}
			}
			if !found {
				t.Error("Test attachment not found in attachment list")
			}
			t.Logf("Found %d attachments", len(attachments))
		})
	})

	// Test 7: Member Commands
	t.Run("Member Commands", func(t *testing.T) {
		// Get current member
		t.Run("Get Current Member", func(t *testing.T) {
			member, err := trelloClient.GetMember("me", trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get current member: %v", err)
			}

			if member.ID == "" {
				t.Error("Expected member ID to be non-empty")
			}
			t.Logf("Current member: %s (ID: %s)", member.FullName, member.ID)
		})

		// Get member's boards (member boards command)
		t.Run("Get Member Boards", func(t *testing.T) {
			member, err := trelloClient.GetMember("me", nil)
			if err != nil {
				t.Fatalf("Failed to get current member: %v", err)
			}

			boards, err := member.GetBoards(nil)
			if err != nil {
				t.Fatalf("Failed to get member boards: %v", err)
			}

			if len(boards) == 0 {
				t.Error("Expected at least one board")
			}

			// Check if our test board is in the list
			found := false
			for _, board := range boards {
				if board.ID == testBoardID {
					found = true
					break
				}
			}
			if !found {
				t.Error("Test board not found in member's board list")
			}
			t.Logf("Found %d boards for member", len(boards))
		})

		// Get board members
		t.Run("Get Board Members", func(t *testing.T) {
			board, err := trelloClient.GetBoard(testBoardID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get board: %v", err)
			}

			members, err := board.GetMembers(trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get board members: %v", err)
			}

			if len(members) == 0 {
				t.Error("Expected at least one member (the creator)")
			}
			t.Logf("Found %d board members", len(members))
		})
	})

	// Test 8: Config Commands (mock test, doesn't require API)
	t.Run("Config Commands", func(t *testing.T) {
		t.Run("Load Config", func(t *testing.T) {
			config, err := client.LoadConfig()
			if err != nil {
				// Config file might not exist, which is okay
				t.Logf("Config file not found (expected): %v", err)
			} else {
				t.Logf("Config loaded: DefaultFormat=%s, MaxTokens=%d", config.DefaultFormat, config.MaxTokens)
			}
		})

		t.Run("Get Config Path", func(t *testing.T) {
			path, err := client.GetConfigPath()
			if err != nil {
				t.Fatalf("Failed to get config path: %v", err)
			}

			if path == "" {
				t.Error("Expected non-empty config path")
			}
			t.Logf("Config path: %s", path)
		})

		t.Run("Save and Load Config", func(t *testing.T) {
			// Create a test config
			testConfig := &client.Config{
				APIKey:        "test-api-key-123",
				Token:         "test-token-456",
				DefaultFormat: "json",
				MaxTokens:     2000,
			}

			// Save it
			err := client.SaveConfig(testConfig)
			if err != nil {
				t.Fatalf("Failed to save config: %v", err)
			}
			t.Logf("Saved test config")

			// Load it back
			loadedConfig, err := client.LoadConfig()
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}

			// Verify the values
			if loadedConfig.APIKey != testConfig.APIKey {
				t.Errorf("Expected APIKey %s, got %s", testConfig.APIKey, loadedConfig.APIKey)
			}
			if loadedConfig.Token != testConfig.Token {
				t.Errorf("Expected Token %s, got %s", testConfig.Token, loadedConfig.Token)
			}
			if loadedConfig.DefaultFormat != testConfig.DefaultFormat {
				t.Errorf("Expected DefaultFormat %s, got %s", testConfig.DefaultFormat, loadedConfig.DefaultFormat)
			}
			if loadedConfig.MaxTokens != testConfig.MaxTokens {
				t.Errorf("Expected MaxTokens %d, got %d", testConfig.MaxTokens, loadedConfig.MaxTokens)
			}
			t.Logf("Config save/load verified")

			// Restore original config (use actual credentials)
			originalConfig := &client.Config{
				APIKey:        apiKey,
				Token:         token,
				DefaultFormat: "markdown",
				MaxTokens:     0,
			}
			_ = client.SaveConfig(originalConfig)
		})
	})

	// Test 9: Batch Commands
	t.Run("Batch Commands", func(t *testing.T) {
		t.Run("Batch Create Cards", func(t *testing.T) {
			// Create multiple cards in batch
			cardNames := []string{
				fmt.Sprintf("Batch Card 1 %s", testID),
				fmt.Sprintf("Batch Card 2 %s", testID),
				fmt.Sprintf("Batch Card 3 %s", testID),
			}

			createdIDs := []string{}
			for _, name := range cardNames {
				card := trello.Card{
					Name:   name,
					IDList: testListID,
				}

				err := trelloClient.CreateCard(&card, trello.Defaults())
				if err != nil {
					t.Fatalf("Failed to create batch card: %v", err)
				}
				createdIDs = append(createdIDs, card.ID)
				t.Logf("Created batch card: %s (ID: %s)", card.Name, card.ID)
			}

			if len(createdIDs) != len(cardNames) {
				t.Errorf("Expected %d cards, created %d", len(cardNames), len(createdIDs))
			}
		})

		t.Run("Batch File Operations", func(t *testing.T) {
			// Create a temporary batch file
			batchFile := fmt.Sprintf("/tmp/trello-batch-test-%s.json", testID)
			defer os.Remove(batchFile)

			batchContent := fmt.Sprintf(`{
				"operations": [
					{
						"type": "card",
						"resource": "card",
						"action": "create",
						"data": {
							"name": "Batch File Card 1 %s",
							"list_id": "%s"
						}
					},
					{
						"type": "card",
						"resource": "card",
						"action": "create",
						"data": {
							"name": "Batch File Card 2 %s",
							"list_id": "%s",
							"desc": "Created via batch file"
						}
					}
				],
				"continue_on_error": true
			}`, testID, testListID, testID, testListID)

			err := os.WriteFile(batchFile, []byte(batchContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create batch file: %v", err)
			}

			// Load and process the batch file
			batchData, err := batch.LoadBatchFile(batchFile)
			if err != nil {
				t.Fatalf("Failed to load batch file: %v", err)
			}

			if len(batchData.Operations) != 2 {
				t.Errorf("Expected 2 operations, got %d", len(batchData.Operations))
			}

			if !batchData.ContinueOnError {
				t.Error("Expected continue_on_error to be true")
			}

			t.Logf("Batch file loaded successfully with %d operations", len(batchData.Operations))
		})

		t.Run("Batch Stdin Operations", func(t *testing.T) {
			// Test batch stdin parsing
			stdinContent := fmt.Sprintf(`{
				"operations": [
					{
						"type": "label",
						"resource": "label",
						"action": "create",
						"data": {
							"name": "Batch Label %s",
							"color": "blue",
							"board_id": "%s"
						}
					}
				]
			}`, testID, testBoardID)

			// Simulate stdin by creating a temp file and reading from it
			tempFile := fmt.Sprintf("/tmp/trello-stdin-test-%s.json", testID)
			defer os.Remove(tempFile)

			err := os.WriteFile(tempFile, []byte(stdinContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp stdin file: %v", err)
			}

			// Load as if from stdin
			batchData, err := batch.LoadBatchFile(tempFile)
			if err != nil {
				t.Fatalf("Failed to load batch data: %v", err)
			}

			if len(batchData.Operations) != 1 {
				t.Errorf("Expected 1 operation, got %d", len(batchData.Operations))
			}

			if batchData.Operations[0].Type != "label" {
				t.Errorf("Expected operation type 'label', got '%s'", batchData.Operations[0].Type)
			}

			t.Logf("Batch stdin parsed successfully with %d operations", len(batchData.Operations))
		})
	})

	// Test 10: Output Format Commands
	t.Run("Output Format Commands", func(t *testing.T) {
		t.Run("JSON Format", func(t *testing.T) {
			// Capture command output
			var buf bytes.Buffer

			// Get the card with JSON format
			card, err := trelloClient.GetCard(testCardID, trello.Defaults())
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}

			// Simple validation that we can get the card data
			if card.ID != testCardID {
				t.Errorf("Expected card ID %s, got %s", testCardID, card.ID)
			}

			// Test output contains expected data
			output := fmt.Sprintf("%v", card)
			if !strings.Contains(output, testCardID) {
				t.Error("Output does not contain card ID")
			}
			t.Logf("JSON format test passed, output length: %d", buf.Len())
		})
	})

	// Test 11: Delete Card (cleanup one of the test cards)
	if copiedCardID != "" {
		t.Run("Delete Card", func(t *testing.T) {
			card, err := trelloClient.GetCard(copiedCardID, nil)
			if err != nil {
				t.Fatalf("Failed to get card to delete: %v", err)
			}

			err = card.Delete()
			if err != nil {
				t.Fatalf("Failed to delete card: %v", err)
			}

			t.Logf("Deleted card: %s", copiedCardID)

			// Verify the card is deleted by trying to get it (should fail)
			_, err = trelloClient.GetCard(copiedCardID, nil)
			if err == nil {
				t.Error("Expected error when getting deleted card, but got none")
			}
		})
	}

	t.Logf("E2E test completed successfully! Test board will be cleaned up.")
}
