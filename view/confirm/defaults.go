package confirm

import "github.com/fatih/color"

const DefaultTitle = "Warning"

var st = color.BlueString("*")
var DefaultPrompt = "                                      " + color.YellowString(";   :   ;") + "\n" +
	"            " + st + "              " + st + "       " + color.YellowString(".   \\_,!,_/   ,") + "\n" +
	"       __            " + st + "              " + color.YellowString("`.,'     `.,'") + "\n" +
	`  ---  \ \______          ` + st + `          ` + color.YellowString(`/         \`) + "\n" +
	"--- " + color.RedString("#") + color.YellowString("#") + color.HiYellowString("#") + "[==______>               " + color.YellowString("~ -- :         : -- ~") + "\n" +
	"  ---  /_/                           " + color.YellowString(`\         /`) + "\n" +
	"                    " + st + "               " + color.YellowString(",'`._   _.'`.") + "\n" +
	"         " + st + "                         " + color.YellowString("'   / `!` \\   `") + "\n" +
	"     " + st + "                     " + st + "          " + color.YellowString(";   :   ;") + "  \n"


