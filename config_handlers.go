package main

import (
	"fmt"
	"log"
	"regexp"
)

func getBlacklist(inService, botName, inServer, inChannel string) (blacklist []string) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						for _, perm := range server.Permissions {
							if perm.BlockList {
								for _, user := range perm.Users {
									blacklist = append(blacklist, user)
								}
							}
						}
					}
				}
			}
		}
	case "irc":
		for _, group := range getChannelGroups(inService, botName, inServer) {
			for _, channel := range group.ChannelIDs {
				if channel == inChannel {
					for _, perm := range group.Permissions {
						if perm.BlockList {
							for _, user := range perm.Users {
								blacklist = append(blacklist, user)
							}
						}
					}
				}
			}
		}
	default:
	}

	return
}

func getChannels(inService, botName, inServer string) (channels []string) {
	Log.Debugf("service: %s, bot: %s, server: %s", inService, botName, inServer)
	switch inService {
	case "discord":
		for bid := range discordGlobal.Bots {
			Log.Debugf("checking for bot: %s", discordGlobal.Bots[bid].BotName)
			if botName == discordGlobal.Bots[bid].BotName {
				Log.Debugf("matched for %s", discordGlobal.Bots[bid].BotName)
				for sid := range discordGlobal.Bots[bid].Servers {
					Log.Debugf("checking for server: %s", discordGlobal.Bots[bid].Servers[sid].ServerID)
					if inServer == discordGlobal.Bots[bid].Servers[sid].ServerID {
						Log.Debugf("matched for %s", discordGlobal.Bots[bid].Servers[sid].ServerID)
						for gid := range discordGlobal.Bots[bid].Servers[sid].ChanGroups {
							Log.Debugf("%s", discordGlobal.Bots[bid].Servers[sid].ChanGroups[gid].ChannelIDs)
							for _, channel := range discordGlobal.Bots[bid].Servers[sid].ChanGroups[gid].ChannelIDs {
								channels = append(channels, channel)
							}
						}
					}
				}
			}
		}
	case "irc":
		for _, bot := range ircGlobal.Bots {
			if bot.BotName == botName {
				for _, group := range bot.ChanGroups {
					for _, channel := range group.ChannelIDs {
						channels = append(channels, channel)
					}
				}
			}
		}
	default:
	}

	return
}

func getChannelGroups(inService, botName, inServer string) (chanGroups []channelGroup) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						chanGroups = server.ChanGroups
					}
				}
			}
		}
	case "irc":
		for _, bot := range ircGlobal.Bots {
			if bot.BotName == botName {
				chanGroups = bot.ChanGroups
			}
		}
	default:
	}

	return
}

func getCommands(inService, botName, inServer, inChannel string) (commands []command) {
	// prep stuff for passing to the parser
	for _, group := range getChannelGroups(inService, botName, inServer) {
		for _, channel := range group.ChannelIDs {
			if inChannel == channel {
				for _, command := range group.Commands {
					commands = append(commands, command)
				}
			}
		}
	}

	return
}

func addCommand(inService, botName, inServer string, channelGroup int) (err error) {

	return
}

func getKeywords(inService, botName, inServer, inChannel string) (keywords []keyword) {
	// prep stuff for passing to the parser
	for _, group := range getChannelGroups(inService, botName, inServer) {
		for _, channel := range group.ChannelIDs {
			if inChannel == channel {
				for _, keyword := range group.Keywords {
					keywords = append(keywords, keyword)
				}
			}
		}
	}

	return
}

func getMentions(inService, botName, inServer, inChannel string) (ping, mention responseArray) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						if inChannel == "DirectMessage" {
							mention = bot.Config.DMResp
						} else {
							for _, group := range server.ChanGroups {
								for _, channel := range group.ChannelIDs {
									if inChannel == channel {
										Log.Debugf("bot was mentioned on channel %s", channel)
										Log.Debugf("ping resp %s", group.Mentions.Ping)
										Log.Debugf("mention resp %s", group.Mentions.Mention)
										ping = group.Mentions.Ping
										mention = group.Mentions.Mention
										return
									}
								}
							}
						}
					}
				}
			}
		}
	case "irc":
		for _, bot := range ircGlobal.Bots {
			if bot.BotName == botName {
				if inChannel == bot.Config.Server.Nickname {
					mention = bot.Config.DMResp
				} else {
					for _, group := range bot.ChanGroups {
						for _, channel := range group.ChannelIDs {
							if inChannel == channel {
								ping = group.Mentions.Ping
								mention = group.Mentions.Mention
								return
							}
						}
					}
				}
			}
		}
	default:
	}

	return
}

func getParsing(inService, botName, inServer, inChannel string) (parseConf parsing) {
	// prep stuff for passing to the parser
	for _, group := range getChannelGroups(inService, botName, inServer) {
		for _, channel := range group.ChannelIDs {
			if inChannel == channel {
				parseConf = group.Parsing
			}
		}
	}

	return
}

func getFilter(inService, botName, inServer string) (filters []filter) {
	// prep stuff for passing to the parser
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						filters = server.Filters
					}
				}
			}
		}
	case "irc":
	default:
	}

	return
}

func getBotParseConfig() (maxLogs int, response, reaction []string, allowIP bool) {
	return botConfig.Parsing.Max, botConfig.Parsing.Response, botConfig.Parsing.Reaction, botConfig.Parsing.AllowIP
}

func getPrefix(inService, botName, inServer string) (prefix string) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						prefix = server.Config.Prefix
					}
				}
			}
		}
	case "irc":
		for _, bot := range ircGlobal.Bots {
			if bot.BotName == botName {
				prefix = bot.Config.Prefix
			}
		}
	default:
	}

	return
}

func getCommandClear(inService, botName, inServer string) (clear bool) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						clear = server.Config.Clear
					}
				}
			}
		}
	default:
	}

	return
}

// getPermissions returns all permissions a user has
func getPermissions(user, inService, botName, inServer string, roles []string) (perms []string) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						for _, serverPerms := range server.Permissions {
							// if user in in the group
							for _, permUser := range serverPerms.Users {
								if user == permUser {
									for _, newPerm := range serverPerms.Permissions {
										perms = append(perms, newPerm)
									}
								}
							}
							// if user has a role
							for _, permRole := range serverPerms.Roles {
								for _, userRole := range roles {
									if userRole == permRole {
										for _, newPerm := range serverPerms.Permissions {
											perms = append(perms, newPerm)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	case "irc":
	default:
	}
	return
}

func hasPerms(user, inService, botName, inServer, requestedPerm string, roles []string) (authorized bool) {
	perms := getPermissions(user, inService, botName, inServer, roles)

	for p := range perms {
		validID, err := regexp.Compile(fmt.Sprintf("%s", perms[p]))
		if err != nil {
			log.Printf("There was an error compiling the regex for the lfg command")
			return
		}
		return validID.MatchString(requestedPerm)
	}

	return
}

func listGroupCommands(inService, botName, inServer string, channelGroup int) (commands []string) {
	serverChannelGroups := getChannelGroups(inService, botName, inServer)

	switch inService {
	case "discord":
		for _, channelCommands := range serverChannelGroups[channelGroup].Commands {
				commands = append(commands, fmt.Sprintf("`%s`",channelCommands.Command))
		}
	case "irc":
	default:
	}
	return
}

func listChannelCommands(inService, botName, inServer, inChannel string) (commands []string) {
	// prep stuff for passing to the parser
	switch inService {
	case "discord":
		for _, group := range getCommands(inService, botName, inServer, inChannel) {
			commands = append(commands, fmt.Sprintf("`%s`",group.Command))
		}
	case "irc":
	default:
	}

	return
}

func listGroupKeywords(inService, botName, inServer string, channelGroup int) (keywords []string) {
	serverChannelGroups := getChannelGroups(inService, botName, inServer)

	switch inService {
	case "discord":
		for _, channelKeywords := range serverChannelGroups[channelGroup].Keywords {
			keywords = append(keywords, fmt.Sprintf("`%s`",channelKeywords.Keyword))
		}
	case "irc":
	default:
	}

	return
}

func listChannelKeywords(inService, botName, inServer, inChannel string) (keywords []string) {
	// prep stuff for passing to the parser
	switch inService {
	case "discord":
		for _, group := range getKeywords(inService, botName, inServer, inChannel) {
			keywords = append(keywords, fmt.Sprintf("`%s`",group.Keyword))
		}
	case "irc":
	default:
	}

	return
}

func listChannelGroups(inService, botName, inServer string) (channelGroups []string) {
	serverChannelGroups := getChannelGroups(inService, botName, inServer)

	switch inService {
	case "discord":
		for gi := range serverChannelGroups {
			newChannel := fmt.Sprintf("group %d:", gi)
			for ci, channel := range serverChannelGroups[gi].ChannelIDs {
				if ci == 0 {
					newChannel = fmt.Sprintf("%s <#%s>", newChannel, channel)
				} else {
					newChannel = fmt.Sprintf("%s, <#%s>", newChannel, channel)
				}
			}
			channelGroups = append(channelGroups, newChannel)
		}
	case "irc":
	default:
	}

	return
}