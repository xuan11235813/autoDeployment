package main

import (
	"bytes"
	"fmt"

	"bufio"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

type nodeItem struct {
	nodeIndex string
	IPaddress string
	userName  string
	password  string
}

type nodeOperationItem struct {
	operationName           string
	operationContent        string
	operationRefinedContent string
	operationPrefix         string
	operationPostfix        string
	isSubstitute            bool
}

var legalOperationName = []string{"copy", "command", "copyN", "getN"}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getRealNameFromPattern(oldString string, index string) string {
	indexPos := strings.Index(oldString, "%")

	if indexPos > -1 {
		preName := oldString[:indexPos]
		afterName := oldString[indexPos+1:]
		refinedName := preName + index + afterName
		return refinedName
	}
	return oldString
}

func main() {
	/* variables */
	var operationLines []string
	var nodeOperations []nodeOperationItem
	var lines []string
	var nodes []nodeItem
	argsWithProg := os.Args
	operationFile := "operation.csv"
	if len(argsWithProg) > 1 {
		operationFile = argsWithProg[1]
		fmt.Printf(operationFile + "\n")
	}
	/* download data files */
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0700)
	}
	/* log file */
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0700)
	}
	dt := time.Now().Unix()
	logFile, err := os.Create("log/log" + strconv.FormatInt(dt, 10) + ".log")
	if err != nil {
		panic("Initialize log file failed.")
	}
	defer logFile.Close()
	/* domain file */
	local, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	localNodeDomainFilePath := local + "/nodeDomain.txt"

	nodeDomainFile, err := os.Open(localNodeDomainFilePath)

	if err != nil {
		fmt.Println("There is no node domain file. Please recheck.")
		panic(err)
	}
	defer nodeDomainFile.Close()

	scanner := bufio.NewScanner(nodeDomainFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	nodeDomainFile.Close()

	localOperationFilePath := local + "/" + operationFile
	nodeOperationFile, err := os.Open(localOperationFilePath)

	if err != nil {
		fmt.Println("There is no operation file. Please recheck")
		panic(err)
	}
	defer nodeOperationFile.Close()

	scanner = bufio.NewScanner(nodeOperationFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		operationLines = append(operationLines, scanner.Text())
	}

	for i, v := range lines {
		nodeString := strings.Split(v, ",")
		if len(nodeString) == 4 {
			var currNodeItem nodeItem
			currNodeItem.nodeIndex = strings.Join(strings.Fields(nodeString[0]), "")
			currNodeItem.IPaddress = strings.Join(strings.Fields(nodeString[1]), "")
			currNodeItem.userName = strings.Join(strings.Fields(nodeString[2]), "")
			currNodeItem.password = strings.Join(strings.Fields(nodeString[3]), "")
			nodes = append(nodes, currNodeItem)
		} else {
			fmt.Printf("wrong input %d, : %s\n", i, v)
		}
	}
	if len(nodes) == 0 {
		panic("no suitable node domain files")
	}

	for i, v := range operationLines {
		operationString := strings.SplitN(v, ",", 2)
		if len(operationString) == 2 {
			var currOperationItem nodeOperationItem
			currOperationItem.operationName = strings.TrimSpace(operationString[0])
			currOperationItem.operationContent = strings.TrimSpace(operationString[1])
			nodeOperations = append(nodeOperations, currOperationItem)
		} else {
			fmt.Printf("wrong input %d, : %s\n", i, v)
			panic("no suitable operations, stop the program")
		}
	}
	if len(nodes) == 0 {
		panic("incorrect operation files, stop the program")
	} else {
		for i, v := range nodeOperations {
			if !stringInSlice(v.operationName, legalOperationName) {
				fmt.Printf("incorrect operation names: %s in line %d\n", v.operationName, i)
				panic("opps!")
			}
			if v.operationName == "copy" {
				localFile := local + v.operationContent
				if _, err := os.Stat(localFile); os.IsNotExist(err) {
					fmt.Printf("No file in current path: %s in line %d\n", v.operationContent, i)
					panic("opps!")
				}
			}

			if v.operationName == "copyN" {
				indexPos := strings.Index(v.operationContent, "%")
				if indexPos > -1 {
					preName := v.operationContent[:indexPos]
					afterName := v.operationContent[indexPos+1:]
					refinedName := preName + afterName
					nodeOperations[i].operationRefinedContent = refinedName
				}
				for _, currNode := range nodes {
					realFile := getRealNameFromPattern(v.operationContent, currNode.nodeIndex)
					localFile := local + realFile
					if _, err := os.Stat(localFile); os.IsNotExist(err) {
						fmt.Printf("No file in current path: %s in line %d\n", localFile, i)
						panic("opps!")
					}
				}

			}

			if v.operationName == "getN" {
				indexPos := strings.Index(v.operationContent, "%")
				if indexPos > -1 {
					preName := v.operationContent[:indexPos]
					afterName := v.operationContent[indexPos+1:]
					nodeOperations[i].operationPrefix = preName
					nodeOperations[i].operationPostfix = afterName
					nodeOperations[i].isSubstitute = true
				} else {
					nodeOperations[i].operationPrefix = ""
					nodeOperations[i].operationPostfix = v.operationContent
					nodeOperations[i].isSubstitute = false
				}
			}
		}
	}

	for _, v := range nodes {
		fmt.Printf("===============Implement node: %s, node index: %s =============\n", v.IPaddress, v.nodeIndex)
		fmt.Fprintf(logFile, "===============Implement node: %s, node index: %s =============\n", v.IPaddress, v.nodeIndex)
		var testOperation bool = true

		for _, opeItem := range nodeOperations {
			if !testOperation {
				break
			}
			fmt.Printf("current operation: %s; with detail: %s\n", opeItem.operationName, opeItem.operationContent)
			fmt.Fprintf(logFile, "current operation: %s; with detail: %s\n", opeItem.operationName, opeItem.operationContent)
			switch opeItem.operationName {
			case "copy":
				{
					localFile := local + opeItem.operationContent
					err := transferFile(v, localFile, "", logFile)
					if err == nil {
						fmt.Printf("Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.operationName, opeItem.operationContent)
						fmt.Fprintf(logFile, "Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.operationName, opeItem.operationContent)
					} else {
						testOperation = false
						fmt.Printf("\n")
						fmt.Fprintf(logFile, "\n")
					}

				}
			case "copyN":
				{
					realFile := getRealNameFromPattern(opeItem.operationContent, v.nodeIndex)
					localFile := local + realFile
					destName := filepath.Base(opeItem.operationRefinedContent)
					err := transferFile(v, localFile, destName, logFile)
					if err == nil {
						fmt.Printf("Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.operationName, opeItem.operationContent)
						fmt.Fprintf(logFile, "Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.operationName, opeItem.operationContent)
					} else {
						testOperation = false
						fmt.Printf("\n")
						fmt.Fprintf(logFile, "\n")
					}

				}
			case "command":
				{
					err := directImplement(v, opeItem.operationContent, logFile)
					if err == nil {
						fmt.Printf("Success in node: %s with operation: %s\n", v.IPaddress, opeItem.operationContent)
						fmt.Fprintf(logFile, "Success in node: %s with operation: %s\n", v.IPaddress, opeItem.operationContent)
					} else {
						testOperation = false
						fmt.Printf("\n")
						fmt.Fprintf(logFile, "\n")
					}
				}
			case "getN":
				{
					targetFileName := ""
					if opeItem.isSubstitute == false {
						targetFileName = opeItem.operationPrefix + opeItem.operationPostfix
					} else {
						targetFileName = opeItem.operationPrefix + v.nodeIndex + opeItem.operationPostfix
					}
					localCopyedFileName := local + "/data" + "/r" + v.nodeIndex + filepath.Base(targetFileName)
					fmt.Printf("targetFile: %s\n localFileName: %s\n", targetFileName, localCopyedFileName)
					err := downloadFile(v, targetFileName, localCopyedFileName, logFile)
					if err == nil {
						fmt.Printf("Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.operationName, opeItem.operationContent)
						fmt.Fprintf(logFile, "Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.operationName, opeItem.operationContent)
					} else {
						fmt.Printf("\n")
						fmt.Fprintf(logFile, "\n")
					}

				}
			}
		}
		fmt.Printf("\n")
		fmt.Fprintf(logFile, "\n")
	}

}
func directImplement(currNode nodeItem, command string, printOutput io.Writer) error {
	sshConfig := &ssh.ClientConfig{
		User: currNode.userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(currNode.password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", currNode.IPaddress+":22", sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to dial: "+err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to create session: "+err.Error())
		return err
	}
	defer session.Close()
	/* excute the command */
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(command)
	if err != nil {
		fmt.Printf("Failed to excute command: " + err.Error() + "\n")
		fmt.Fprintf(printOutput, "Failed to excute command: "+err.Error()+"\n")
		return err
	}

	fmt.Printf("Output: " + stdoutBuf.String() + "\n")
	fmt.Fprintf(printOutput, "Output: "+stdoutBuf.String()+"\n")

	return nil
}
func transferFile(currNode nodeItem, filePath string, destName string, printOutput io.Writer) error {
	sshConfig := &ssh.ClientConfig{
		User: currNode.userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(currNode.password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", currNode.IPaddress+":22", sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to dial: "+err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to create session: "+err.Error())
		return err
	}
	defer session.Close()

	dest := "/radar/" + destName
	err = scp.CopyPath(filePath, dest, session)
	if err != nil {
		fmt.Printf("Transfering file error with the destination: " + err.Error())
		fmt.Fprintf(printOutput, "Transfering file error with the destination: "+err.Error())
		return err
	}
	return nil
}

func downloadFile(currNode nodeItem, remotePath string, destLocalPath string, printOutput io.Writer) error {
	sshConfig := &ssh.ClientConfig{
		User: currNode.userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(currNode.password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", currNode.IPaddress+":22", sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to dial: "+err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to create session: "+err.Error())
		return err
	}
	defer session.Close()
	r, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(destLocalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd := "cat " + remotePath
	if err := session.Start(cmd); err != nil {
		return err
	}

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}
