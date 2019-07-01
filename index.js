#!/usr/bin/env node
require("dotenv").config();
const axios = require("axios");

// Store Slack access token in a .env file in the root folder
const accessToken = process.env.SLACK_ACCESS_TOKEN;

async function getConversationId(conversation) {

  // Trim the pound sign if it was provided in the conversation argument
  if (conversation.startsWith('#')) {
    conversation = conversation.slice(1)
  }

  // Destructuring to pull out the channels result from the response
  const { data: { channels } } = await axios.get(`https://slack.com/api/conversations.list?token=${accessToken}&pretty=1`)

  // Find the channel with matching conversation name and pull out its ID
  const conversationId = channels.find(chan => chan.name === conversation).id

  return conversationId
}

async function getConversationMembers(conversationId) {

  async function queryForMembers(membersCollection = [], nextCursor = '') {
    let allConversationMembers = [];

    // Base case: members exist, no next cursor, we've retrieved all the members
    if (membersCollection.length > 0 && !nextCursor) {
      return membersCollection;
    }

    try {
      const response = await axios.get(
        `https://slack.com/api/conversations.members?token=${accessToken}&channel=${conversationId}${nextCursor}`
      );

      // Looks ugly, but pulling out members and next_cursor
      let {
        data: {
          members,
          response_metadata: { next_cursor }
        }
      } = response;

      // Cursor seems to always end with =
      // Replace with URL encoding, %3D
      if (next_cursor) {
        next_cursor = `&cursor=${next_cursor.replace('=', '%3D')}`
      }

      // Combine member grouping with previous
      allConversationMembers = [...members, ...membersCollection]

      return await queryForMembers(allConversationMembers, next_cursor)
    }
    catch (err) {
      console.log(err);
    }

  }

  // Kick off recursive call
  return await queryForMembers()
}

async function inviteMembers(sourceChannel, targetChannel) {

  // Basic check for 2 arguments
  if (process.argv.length !== 4) {
    return console.log("Please supply a source channel and a target channel")
  }

  // Retrieve the ID of the source and target channels for member access
  const sourceId = await getConversationId(sourceChannel)
  const targetId = await getConversationId(targetChannel)

  // Retrieve member lists for comparison
  const sourceMembers = await getConversationMembers(sourceId)
  const targetMembers = await getConversationMembers(targetId)

  // Combine member lists and remove duplicate members.
  // Filter out members that are already in target channel.
  // This prevents cants_invite_self and already_in_channel errors
  const membersToInvite = [...new Set([...sourceMembers, ...targetMembers])]
    .filter(member => !targetMembers.includes(member))

  // Slack API limits the number of invites.
  // Make a new request 20 entries at a time while members remain in the array.
  while (membersToInvite.length > 0) {
    const batch = membersToInvite.splice(0, 20);

    // result currently unused--possible future use.
    const result = await axios.post(`https://slack.com/api/conversations.invite?token=${accessToken}&channel=${targetId}&users=${batch.join('%2C')}`)
    console.log(`${batch.length} members invited.`)
  }
}

inviteMembers(process.argv[2], process.argv[3])
