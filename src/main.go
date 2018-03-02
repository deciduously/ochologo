package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	h "github.com/deciduously/helpers"
	pc "github.com/deciduously/playingcards"
)

type data struct {
	compHand, playerHand, deck, discard pc.Deck
	gameState                           string
	activeSuit                          pc.Suit
}

func main() {

	fmt.Println("Welcome to Ocho LoGo!  When you run out of cards, you WIN!")
	fmt.Println("This is a two-player variant, you and the computer are dealt 7 cards each to begin.")

	d := data{}

	initGame(&d)

	for d.gameState == "open" {
		playerTurn(&d)
		if d.playerHand.Empty() {
			d.gameState = "win"
			break
		}
		compTurn(&d)
		if d.compHand.Empty() {
			d.gameState = "loss"
		}
	}
	if d.gameState == "win" {
		fmt.Println("You win!")
	} else {
		fmt.Println("You lose!")
	}

}

func compTurn(d *data) {
	fmt.Println("\n--------COMPUTER--------")
	suit := d.activeSuit
	val := d.discard.Peek().Value
	canPlay := false
	play := pc.Card{}
	//once through for suit || val match
	for _, c := range d.compHand.Cards {
		if c.Value == val || c.Suit == suit {
			play = c
			canPlay = true
			break
		}
	}
	//...and once through afterwards for wilds, ensuring it's a last resort
	for  _,c := range d.compHand.Cards {
		if c.Value == 0x8 {
			play = c
			canPlay = true
			break
		}
	}
	if canPlay {
		e := playCard(d, "c", play)
		h.CheckErr(e)
		fmt.Printf("Computer played the %v\n", play.Readable())
	} else {
		drawCount := 0
		for {
			draw, err := d.deck.Draw(1)
			h.CheckErr(err)
			if d.deck.Empty() {
				flipDeck(d)
			}
			d.compHand.Push(draw.Cards[0])
			drawCount++
			if draw.Cards[0].Value == val || draw.Cards[0].Suit == suit || draw.Cards[0].Value == 0x8 {
				fmt.Printf("Computer drew %d and played the %v\n", drawCount, draw.Cards[0].Readable())
				e := playCard(d, "c", draw.Cards[0])
				h.CheckErr(e)
				break
			}
		}
	}
	if d.compHand.Len() == 1 {
		fmt.Println("Computer has ONE CARD LEFT!")
	}
}

func flipDeck(d *data) {
	fmt.Println("Flipping discard and reshuffling deck...")
	lastCard := d.discard.Pop()
	for i := 0; i < d.discard.Len(); i++ {
		d.deck.Push(d.discard.Pop())
	}
	d.deck.Shuffle()
	d.discard.Push(lastCard)
}

func initGame(d *data) {
	d.deck = pc.NewDeck()
	d.deck.Shuffle()
	d.gameState = "open"
	var err error

	d.playerHand, err = d.deck.Draw(7)
	h.CheckErr(err)

	d.compHand, err = d.deck.Draw(7)
	h.CheckErr(err)

	d.discard, err = d.deck.Draw(1)
	h.CheckErr(err)

	d.activeSuit = d.discard.Peek().Suit
}

//playCard accepts a pointer to the data struct d, either "c" for computer or "p" for the player,
//and a Card to remove from their hand and add to the discard
func playCard(d *data, player string, c pc.Card) error {
	switch player {
	case "c":
		e := d.compHand.Remove(c)
		h.CheckErr(e)
		d.discard.Push(c)
		if c.Value == 0x2 {
			d.activeSuit = c.Suit
			if d.deck.Len() <= 2 {
				flipDeck(d)
			}
			fmt.Println("It's a two!  You draw two cards")
			n, e := d.deck.Draw(2)
			h.CheckErr(e)
			for _, v := range n.Cards {
				d.playerHand.Push(v)
			}
		} else if c.Value == 0x8 {
			compCounts := make(map[pc.Suit]int)
			for _, c := range d.compHand.Cards {
				switch c.Suit {
				case pc.Hearts:
					compCounts[pc.Hearts]++
				case pc.Diamonds:
					compCounts[pc.Diamonds]++
				case pc.Clubs:
					compCounts[pc.Clubs]++
				case pc.Spades:
					compCounts[pc.Spades]++
				}
			}
			var wants pc.Suit
			for k := range compCounts {
				if compCounts[k] > compCounts[wants] {
					wants = k
				}
			}
			switch wants {
			case pc.Hearts:
				fmt.Println("Computer requests Hearts!")
				d.activeSuit = pc.Hearts
			case pc.Diamonds:
				fmt.Println("Computer requests Diamonds!")
				d.activeSuit = pc.Diamonds
			case pc.Clubs:
				fmt.Println("Computer requests Clubs!")
				d.activeSuit = pc.Clubs
			case pc.Spades:
				fmt.Println("Computer requests Spades!")
				d.activeSuit = pc.Spades
			}
		} else {
			d.activeSuit = c.Suit
		}
	case "p":
		e := d.playerHand.Remove(c)
		h.CheckErr(e)
		d.discard.Push(c)
		if c.Value == 0x2 {
			d.activeSuit = c.Suit
			if d.deck.Len() <= 2 {
				flipDeck(d)
			}
			fmt.Println("It's a two!  Computer draws two cards")
			n, e := d.deck.Draw(2)
			h.CheckErr(e)
			for _, v := range n.Cards {
				d.compHand.Push(v)
			}
		} else if c.Value == 0x8 {
		ploco:
			for {
				var input string
				fmt.Println("Ocho LoGo pide... (h)earts, (d)iamonds, (c)lubs, (s)pades: ")
				_, e := fmt.Scanf("%s", &input)
				if e != nil {
					continue
				}
				switch input {
				case "h":
					d.activeSuit = pc.Hearts
					fmt.Println("Player requests Hearts!")
					break ploco
				case "d":
					d.activeSuit = pc.Diamonds
					fmt.Println("Player requests Diamonds!")
					break ploco
				case "c":
					d.activeSuit = pc.Clubs
					fmt.Println("Player requests Clubs")
					break ploco
				case "s":
					d.activeSuit = pc.Spades
					fmt.Println("Player requests Spades")
					break ploco
				default:
					fmt.Println("What was that?")
					continue
				}
			}
		} else {
			d.activeSuit = c.Suit
		}
	default:
		return fmt.Errorf("Player string MUST be \"c\" or \"p\", I was passed %s", player)
	}
	return nil
}

func playerTurn(d *data) {
	fmt.Println("\n----------YOU-----------")
	var canPlay bool
	var hasWild bool
	suit := d.activeSuit
	val := d.discard.Peek().Value
	fmt.Println("\nCurrent discard: " + d.discard.Peek().Readable() + "\n\nYour Hand:")

	for i, c := range d.playerHand.Cards {
		if c.Value == val || c.Suit == suit {
			canPlay = true
		} else if c.Value == 0x8 {
			hasWild = true
		}
		fmt.Println(strconv.Itoa(i+1) + " - " + c.Readable())
	}
	if canPlay || hasWild {
		var ps int
		var choice pc.Card
		for {
			fmt.Print("\nYour choice: ")
			ps = h.GetInt()
			choice = d.playerHand.Cards[ps-1]
			if ps > 0 && ps <= d.playerHand.Len() && (choice.Value == val || choice.Suit == suit || (choice.Value == 0x8 && !canPlay)) {
				break
			} else if choice.Value == 0x8 && canPlay {
				fmt.Println("You can only use your 8 as a last resort!")
			} else {
				fmt.Println("Invalid!  Try again\nCurrent discard: " + d.discard.Peek().Readable())
			}
		}
		play := choice
		e := playCard(d, "p", play)
		h.CheckErr(e)
		fmt.Printf("You play the %v\n", play.Readable())
	} else {
		for {
			fmt.Println("Press enter to draw a card")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			draw, err := d.deck.Draw(1)
			h.CheckErr(err)
			if d.deck.Empty() {
				flipDeck(d)
			}
			d.playerHand.Push(draw.Cards[0])
			if draw.Cards[0].Value == val || draw.Cards[0].Suit == suit || draw.Cards[0].Value == 0x8 {
				fmt.Printf("You draw and play the %v\n", draw.Cards[0].Readable())
				e := playCard(d, "p", draw.Cards[0])
				h.CheckErr(e)
				break
			} else {
				fmt.Printf("You can't play the %v\n", draw.Cards[0].Readable())
			}
		}
	}
	if d.playerHand.Len() == 1 {
		fmt.Println("You have ONE CARD LEFT!")
	}
}
