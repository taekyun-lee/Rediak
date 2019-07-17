package KangDB

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main(){
	scanner := bufio.NewReader(os.Stdin)
	fmt.Println("Scannertest!!!!! ")
	for{
		fmt.Print(">>")
		text, _ := scanner.ReadString('\n')
		listcmd := strings.Split(text," ")
		cmd := listcmd[0]
		switch cmd{
		case "GET":
			fmt.Println("getman : ",listcmd)
		case "SET":
			fmt.Println("getman : ",listcmd)
		case "EXIT":
			break
		default:
			fmt.Println("Fail cmd")
		}


	}
}
