# Streamer
Streamer is used to create data streams that leverage Golang channels.
 
 ## Why 🤔 Though
So hear's the thing I love channels I think they are on of the most powerful tools in golang.
But there are some things about channels that I think can increase their utility.

1. You don't know if a channel is closed without receiving from it.
	- If you send to a closed channel you panic! But there is no great way to know if someone closed it before sending.
2.  You can not close a closed channel.
	- If multiple people want to close a channel because of an error at the same time you have a race condition.
3. You can not guarantee that a send to channel will not block.
	- When dealing with a high volume of information that sends to a channel. You may not care about older values not processed in a timely manner. For this instance it would be nice if the channel worked as first in first out data structure
4. Zero value of a channel is nil. 
	- If you attempt to use the channel without initializing it then you will panic or block forever.
 
There is an argument to be made that if you desire some of these behaviors you can redesign you program to avoid the scenarios. Though this is probably true that does not make the new design inherently better or easier to read and maintain. 