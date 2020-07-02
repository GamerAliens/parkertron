package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func discordCommandHandler(message, user, inService, botName, inServer, inChannel string, roles []string) (response, reaction []string, err error) {
	messageLines := strings.Split(message, "\n")
	stringMatch, err := regexp.Compile(`(add|list|mod|del)(?:\s+|$)(changroups?|keywords?|commands?)(?:\s+|$)(?:(\d+|page \d+|<#\d{18}>|))(?:\s?|$)(?:(\d+|page \d+|exact|)|)(?:\s+|$)(?:(.*)|)`)
	if err != nil {
		return
	}

	matchedString := stringMatch.FindStringSubmatch(messageLines[0])

	if len(matchedString) == 0 {
		return
	}

	action := matchedString[1]
	option := strings.TrimSuffix(matchedString[2], "s")
	page := 1
	var group int
	var channel string
	var newExact bool
	newKeyword := fmt.Sprintf("")
	newResponse := messageLines[1:]

	if !hasPerms(user, inService, botName, inServer, fmt.Sprintf("%s.%s", option, action), roles) {
		return
	}

	if strings.HasPrefix(matchedString[3], "<#") && strings.HasSuffix(matchedString[3], ">") {
		channel = strings.TrimSuffix(strings.TrimPrefix(matchedString[3], "<#"), ">")
	} else if matchedString[3] != "" {
		if len(strings.Split(matchedString[3], " ")) > 1 && strings.Split(matchedString[3], " ")[0] == "page" {
			if page, err = strconv.Atoi(strings.Split(matchedString[3], " ")[1]); err != nil {
				return
			}
			channel = inChannel
		} else if len(strings.Split(matchedString[3], " ")) == 1 {
			if group, err = strconv.Atoi(matchedString[3]); err != nil {
				return
			}
		}
	} else {
		channel = inChannel
	}

	if matchedString[4] == "exact" {
		newExact = true
	} else if matchedString[4] != "" {
		if len(strings.Split(matchedString[4], " ")) == 1 {
			if page, err = strconv.Atoi(matchedString[4]); err != nil {
				return
			}
		} else if len(strings.Split(matchedString[4], " ")) > 1 && strings.Split(matchedString[4], " ")[0] == "page" {
			if page, err = strconv.Atoi(strings.Split(matchedString[4], " ")[1]); err != nil {
				return
			}
		}
	}

	if action == "list" {
		response = []string{fmt.Sprintf("")}
		if response, reaction, err = readMatch(option, inService, botName, inServer, channel, group, page); err != nil {
			return
		}

	}

	if action == "add" {
		if response, reaction, err = creatMatch(option, inService, botName, inServer, newKeyword, newResponse, group, newExact); err != nil {
			return
		}
	} else if action == "mod" {

	} else if action == "del" {

	}

	return
}

// TODO: re-write CRUD handlers

func creatMatch(option, inService, botName, inServer, newKeyword string, newResponse []string, groupID int, exact bool) (response, reaction []string, err error) {

	return
}

func readMatch(option, inService, botName, inServer, inChannel string, groupID, page int) (response, reaction []string, err error) {
	// if no channel was passed along it's for a group
	if inChannel == "" && option != "changroup"{
		// group
		response = []string{fmt.Sprintf("listing %ss for group %d", option, groupID)}
	} else if option != "changroup"{
		// channel
		response = []string{fmt.Sprintf("listing %ss for channel <#%s>", option, inChannel)}
	} else {
		response = []string{fmt.Sprintf("listing channel groups")}
	}

	var listArray []string

	switch option {
	case "keyword":
		if inChannel == "" {
			listArray = listGroupKeywords(inService, botName, inServer, groupID)
		} else {
			listArray = listChannelKeywords(inService, botName, inServer, inChannel)
		}

	case "command":
		if inChannel == "" {
			listArray = listGroupCommands(inService, botName, inServer, groupID)
		} else {
			listArray = listChannelCommands(inService, botName, inServer, inChannel)
		}
	case "changroup":
		listArray = listChannelGroups(inService, botName, inServer)
	}

	var newList string

	Log.Debug(listArray)


	start := (page-1) * 20

	if option == "changroup" {
		newList = strings.Join(listArray, "\n")
	} else {
		for count, listItem := range listArray {
			if start == count {
				newList = fmt.Sprintf("%s", listItem)
			} else if start < count && count < start + 20 {
				newList = fmt.Sprintf("%s, %s", newList, listItem)
			}
		}
	}

	totalPages := 1

	if len(listArray)%20 > 0 {
		totalPages = (len(listArray)/20) + 1
	}

	response = append(response, newList)
	if totalPages > 1 {
		response = append(response, fmt.Sprintf("this was page %d of %d to show more pages `list %ss page #`", page, totalPages, option))
	}
	
	return
}

func updateMatch(option, inService, botName, inServer string, groupID int) (response, reaction []string, err error) {
	return
}

func deleteMatch(option, inService, botName, inServer string, groupID int) (response, reaction []string, err error) {
	return
}