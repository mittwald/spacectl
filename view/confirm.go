package view

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

var st = color.BlueString("*")
var prompt = "                                      " + color.YellowString(";   :   ;") + "\n" +
	"            " + st + "              " + st + "       " + color.YellowString(".   \\_,!,_/   ,") + "\n" +
	"       __            " + st + "              " + color.YellowString("`.,'     `.,'") + "\n" +
	`  ---  \ \______          ` + st + `          ` + color.YellowString(`/         \`) + "\n" +
	"--- " + color.RedString("#") + color.YellowString("#") + color.HiYellowString("#") + "[==______>               " + color.YellowString("~ -- :         : -- ~") + "\n" +
	"  ---  /_/                           " + color.YellowString(`\         /`) + "\n" +
	"                    " + st + "               " + color.YellowString(",'`._   _.'`.") + "\n" +
	"         " + st + "                         " + color.YellowString("'   / `!` \\   `") + "\n" +
	"     " + st + "                     " + st + "          " + color.YellowString(";   :   ;") + "  \n"

func Confirm(msg string, body string) (bool, error) {
	if msg == "" {
		msg = "This action is destructive!"
	}

	fmt.Fprintln(color.Output, prompt)

	color.Red("WARNING\n")
	color.Red("  %s", msg)
	color.Red("  Are you ABSOLUTELY sure you want to continue?")
	fmt.Println("")

	if body != "" {
		fmt.Println(body)
		fmt.Println("")
	}

	for {
		fmt.Fprint(color.Output, "Enter "+color.YellowString("y")+" or "+color.YellowString("n")+": ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		input = strings.TrimRight(input, "\n")

		if input == "y" {
			return true, nil
		} else if input == "n" {
			return false, nil
		}
	}
}
