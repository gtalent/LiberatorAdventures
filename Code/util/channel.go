/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package util

//Used to route messages to the main channel.
type ChannelLine struct {
	name    string
	channel chan string
}

//Creates a new ChannelLine to route to the given channel.
func NewChannelLine(name string, channel chan string) *ChannelLine {
	line := new(ChannelLine)
	line.channel = channel
	line.name = name
	return line
}

//Forwards the given message to this channel, specifying the name of the channel as the source of the message.
func (me *ChannelLine) Put(message string) {
	me.channel <- (me.name + ":\t" + message)
}

