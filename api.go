package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type (
	// User represents a Telegram user or bot.
	User struct {
		// Unique identifier for this user or bot.
		ID int `json:"id"`
		// User‘s or bot’s first name.
		FirstName string `json:"first_name"`
		// Optional. User‘s or bot’s last name.
		LastName string `json:"last_name"`
		// Optional. User‘s or bot’s username.
		Username string `json:"username"`
	}

	// Chat represents a chat.
	Chat struct {
		// Unique identifier for this chat. This number may be greater than
		// 32 bits.
		ID int64 `json:"id"`
		// Type of chat, can be either “private”, “group”, “supergroup” or
		// “channel”.
		Type string `json:"type"`
		// Optional. Title, for supergroups, channels and group chats.
		Title string `json:"title"`
		// Optional. Username, for private chats, supergroups and channels if
		// available.
		Username string `json:"username"`
		// Optional. First name of the other party in a private chat.
		FirstName string `json:"first_name"`
		// Optional. Last name of the other party in a private chat.
		LastName string `json:"last_name"`
		// Optional. True if a group has ‘All Members Are Admins’ enabled.
		AllMembersAreAdministrators bool `json:"all_members_are_administrators"`
	}

	// Message is abstract type of telegram incoming messages.
	Message struct {
		// Unique message identifier.
		MessageID int `json:"message_id"`
		// Optional. Sender, can be empty for messages sent to channels.
		From *User `json:"from"`
		// Date the message was sent in Unix time.
		Date uint64 `json:"date"`
		// Conversation the message belongs to.
		Chat *Chat `json:"chat"`
		// Optional. For forwarded messages, sender of the original message.
		ForwardFrom *User `json:"forward_from"`
		// Optional. For messages forwarded from a channel, information about
		// the original channel.
		ForwardFromChat *Chat `json:"forward_from_chat"`
		// Optional. For forwarded messages, date the original message was sent
		// in Unix time.
		ForwardDate uint64 `json:"forward_date"`
		// Optional. For replies, the original message. Note that the Message
		// object in this field will not contain further reply_to_message fields
		// even if it itself is a reply.
		ReplyToMessage *Message `json:"reply_to_message"`
		// Optional. Date the message was last edited in Unix time.
		EditDate uint64 `json:"edit_date"`
		// Optional. For text messages, the actual UTF-8 text of the message,
		// 0-4096 characters.
		Text string `json:"text"`
		// Optional. For text messages, special entities like usernames, URLs,
		// bot commands, etc. that appear in the text.
		Entities []MessageEntity `json:"entities"`
		// Optional. Message is an audio file, information about the file.
		Audio *Audio `json:"audio"`
		// Optional. Message is a general file, information about the file.
		Document *Document `json:"document"`
		// Optional. Message is a game, information about the game.
		Game *Game `json:"game"`
		// Optional. Message is a photo, available sizes of the photo.
		Photo []PhotoSize `json:"photo"`
		// Optional. Message is a sticker, information about the sticker.
		Sticker *Sticker `json:"sticker"`
		// Optional. Message is a video, information about the video.
		Video *Video `json:"video"`
		// Optional. Message is a voice message, information about the file.
		Voice *Voice `json:"voice"`
		// Optional. Caption for the document, photo or video, 0-200 characters.
		Caption string `json:"caption"`
		// Optional. Message is a shared contact, information about the contact.
		Contact *Contact `json:"contact"`
		// Optional. Message is a shared location, information about the
		// location.
		Location *Location `json:"location"`
		// Optional. Message is a venue, information about the venue.
		Venue *Venue `json:"venue"`
		// Optional. A new member was added to the group, information about them
		// (this member may be the bot itself).
		NewChatMember *User `json:"new_chat_member"`
		// Optional. A member was removed from the group, information about them
		// (this member may be the bot itself).
		LeftChatMember *User `json:"left_chat_member"`
		// Optional. A chat title was changed to this value.
		NewChatTitle string `json:"new_chat_title"`
		// Optional. A chat photo was change to this value.
		NewChatPhoto []PhotoSize `json:"new_chat_photo"`
		// Optional. Service message: the chat photo was deleted.
		DeleteChatPhoto bool `json:"delete_chat_photo"`
		// Optional. Service message: the group has been created
		GroupChatCreated bool `json:"group_chat_created"`
		// Optional. Service message: the supergroup has been created. This
		// field can‘t be received in a message coming through updates, because
		// bot can’t be a member of a supergroup when it is created. It can only
		// be found in reply_to_message if someone replies to a very first
		// message in a directly created supergroup.
		SupergroupChatCreated bool `json:"supergroup_chat_created"`
		// Optional. Service message: the channel has been created. This field
		// can‘t be received in a message coming through updates, because bot
		// can’t be a member of a channel when it is created. It can only be
		// found in reply_to_message if someone replies to a very first message
		// in a channel.
		ChannelChatCreated bool `json:"channel_chat_created"`
		// Optional. The group has been migrated to a supergroup with the
		// specified identifier. This number may be greater than 32 bits.
		MigrateToChatID int64 `json:"migrate_to_chat_id"`
		// Optional. The supergroup has been migrated from a group with the
		// specified identifier. This number may be greater than 32 bits.
		MigrateFromChatID int64 `json:"migrate_from_chat_id"`
		// Optional. Specified message was pinned. Note that the Message object
		// in this field will not contain further reply_to_message fields even
		// if it is itself a reply.
		PinnedMessage *Message `json:"pinned_message"`
	}

	// Update represents an incoming update. Only one of the optional parameters
	// can be present in any given update.
	Update struct {
		// The update‘s unique identifier. Update identifiers start from a
		// certain positive number and increase sequentially. This ID becomes
		// especially handy if you’re using Webhooks, since it allows you to
		// ignore repeated updates or to restore the correct update sequence,
		// should they get out of order.
		UpdateID int `json:"update_id"`
		// Optional. New incoming message of any kind — text, photo, sticker,
		// etc..
		Message *Message `json:"message"`
		// Optional. New version of a message that is known to the bot and was
		// edited.
		EditedMessage *Message `json:"edited_message"`
		// Optional. New incoming inline query.
		InlineQuery *InlineQuery `json:"inline_query"`
		// Optional. The result of an inline query that was chosen by a user and
		// sent to their chat partner..
		ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
		// Optional. New incoming callback query.
		CallbackQuery *CallbackQuery `json:"callback_query"`
	}

	// MessageEntity represents one special entity in a text message. For
	// example, hashtags, usernames, URLs, etc.
	MessageEntity struct {
		// Type of the entity. Can be mention (@username), hashtag, bot_command,
		// url, email, bold (bold text), italic (italic text), code (monowidth
		// string), pre (monowidth block), text_link (for clickable text URLs),
		// text_mention (for users without usernames).
		Type string `json:"type"`
		// Offset in UTF-16 code units to the start of the entity.
		Offset int `json:"offset"`
		// Length of the entity in UTF-16 code units.
		Length int `json:"length"`
		// Optional. For “text_link” only, url that will be opened after user
		// taps on the text.
		URL string `json:"url"`
		// Optional. For “text_mention” only, the mentioned user.
		User *User `json:"user"`
	}

	Game struct {
		// Title of the game.
		Title string `json:"title"`
		// Description of the game.
		Description string `json:"description"`
		// Photo that will be displayed in the game message in chats.
		Photo []PhotoSize `json:"photo"`
		// Optional. Brief description of the game or high scores included in
		// the game message. Can be automatically edited to include current high
		// scores for the game when the bot calls setGameScore, or manually
		// edited using editMessageText. 0-4096 characters.
		Text string `json:"text"`
		// Optional. Special entities that appear in text, such as usernames,
		// URLs, bot commands, etc.
		TextEntities []MessageEntity `json:"text_entities"`
		// Optional. Animation that will be displayed in the game message in
		// chats. Upload via BotFather.
		Animation *Animation `json:"animation"`
	}

	Animation struct {
		// Unique file identifier.
		FileID string `json:"file_id"`
		// Optional. Animation thumbnail as defined by sender.
		Thumb *PhotoSize `json:"thumb"`
		// Optional. Original animation filename as defined by sender.
		FileName string `json:"file_name"`
		// Optional. MIME type of the file as defined by sender.
		MimeType string `json:"mime_type"`
		// Optional. File size.
		FileSize int `json:"file_size"`
	}

	// PhotoSize represents one size of a photo or a file / sticker thumbnail.
	PhotoSize struct {
		// Unique identifier for this file.
		FileID string `json:"file_id"`
		// Photo width.
		Width int `json:"width"`
		// Photo height.
		Height int `json:"height"`
		// Optional. File size.
		FileSize int `json:"file_size"`
	}

	// Audio represents an audio file to be treated as music by the Telegram
	// clients.
	Audio struct {
		// Unique identifier for this file.
		FileID string `json:"file_id"`
		// Duration of the audio in seconds as defined by sender.
		Duration int `json:"duration"`
		// Optional. Performer of the audio as defined by sender or by audio
		// tags.
		Performer string `json:"performer"`
		// Optional. Title of the audio as defined by sender or by audio tags.
		Title string `json:"title"`
		// Optional. MIME type of the file as defined by sender.
		MimeType string `json:"mime_type"`
		// Optional. File size.
		FileSize int `json:"file_size"`
	}

	// Document represents a general file (as opposed to photos, voice messages
	// and audio files).
	Document struct {
		// Unique file identifier.
		FileID string `json:"file_id"`
		// Optional. Document thumbnail as defined by sender.
		Thumb *PhotoSize `json:"thumb"`
		// Optional. Original filename as defined by sender.
		FileName string `json:"file_name"`
		// Optional. MIME type of the file as defined by sender.
		MimeType string `json:"mime_type"`
		// Optional. File size.
		FileSize int `json:"file_size"`
	}

	// Sticker represents a sticker.
	Sticker struct {
		// Unique identifier for this file.
		FileID string `json:"file_id"`
		// Sticker width.
		Width int `json:"width"`
		// Sticker height.
		Height int `json:"height"`
		// Optional. Sticker thumbnail in .webp or .jpg format.
		Thumb *PhotoSize `json:"thumb"`
		// Optional. Emoji associated with the sticker.
		Emoji string `json:"emoji"`
		// Optional. File size.
		FileSize int `json:"file_size"`
	}

	// Video represents a video file.
	Video struct {
		// Unique identifier for this file.
		FileID string `json:"file_id"`
		// Video width as defined by sender.
		Width int `json:"width"`
		// Video height as defined by sender.
		Height int `json:"height"`
		// Duration of the video in seconds as defined by sender.
		Duration int `json:"duration"`
		// Optional. Video thumbnail.
		Thumb *PhotoSize `json:"thumb"`
		// Optional. Mime type of a file as defined by sender.
		MimeType string `json:"mime_type"`
		// Optional. File size.
		FileSize int `json:"file_size"`
	}

	// Voice represents a voice note.
	Voice struct {
		// Unique identifier for this file.
		FileID string `json:"file_id"`
		// Duration of the audio in seconds as defined by sender.
		Duration int `json:"duration"`
		// Optional. MIME type of the file as defined by sender.
		MimeType string `json:"mime_type"`
		// Optional. File size.
		FileSize int `json:"file_size"`
	}

	// Contact represents a phone contact.
	Contact struct {
		// Contact's phone number.
		PhoneNumber string `json:"phone_number"`
		// Contact's first name.
		FirstName string `json:"first_name"`
		// Optional. Contact's last name.
		LastName string `json:"last_name"`
		// Optional. Contact's user identifier in Telegram.
		UserID int `json:"user_id"`
	}

	// Location represents a point on the map.
	Location struct {
		// Longitude as defined by sender.
		Longitude float64 `json:"longitude"`
		// Latitude as defined by sender.
		Latitude float64 `json:"latitude"`
	}

	// Venue represents a venue.
	Venue struct {
		// Venue location.
		Location *Location `json:"location"`
		// Name of the venue.
		Title string `json:"title"`
		// Address of the venue.
		Address string `json:"address"`
		// Optional. Foursquare identifier of the venue.
		FoursquareID string `json:"foursquare_id"`
	}

	// InlineQuery represents an incoming inline query. When the user sends an
	// empty query, your bot could return some default or trending results.
	InlineQuery struct {
		// Unique identifier for this query.
		ID string `json:"id"`
		// Sender.
		From *User `json:"from"`
		// Optional. Sender location, only for bots that request user location.
		Location *Location `json:"location"`
		// Text of the query (up to 512 characters).
		Query string `json:"query"`
		// Offset of the results to be returned, can be controlled by the bot.
		Offset string `json:"offset"`
	}

	// ChosenInlineResult represents a result of an inline query that was chosen
	// by the user and sent to their chat partner.
	ChosenInlineResult struct {
		// The unique identifier for the result that was chosen.
		ResultID string `json:"result_id"`
		// The user that chose the result.
		From *User `json:"from"`
		// Optional. Sender location, only for bots that require user location.
		Location *Location `json:"location"`
		// Optional. Identifier of the sent inline message. Available only if
		// there is an inline keyboard attached to the message. Will be also
		// received in callback queries and can be used to edit the message.
		InlineMessageID string `json:"inline_message_id"`
		// The query that was used to obtain the result.
		Query string `json:"query"`
	}

	// CallbackQuery represents an incoming callback query from a callback
	// button in an inline keyboard. If the button that originated the query was
	// attached to a message sent by the bot, the field message will be present.
	// If the button was attached to a message sent via the bot (in inline
	// mode), the field inline_message_id will be present. Exactly one of the
	// fields data or game_short_name will be present.
	CallbackQuery struct {
		// Unique identifier for this query.
		ID string `json:"id"`
		// Sender.
		From *User `json:"from"`
		// Optional. Message with the callback button that originated the query.
		// Note that message content and message date will not be available if
		// the message is too old.
		Message *Message `json:"message"`
		// Optional. Identifier of the message sent via the bot in inline mode,
		// that originated the query..
		InlineMessageID string `json:"inline_message_id"`
		// Identifier, uniquely corresponding to the chat to which the message
		// with the callback button was sent. Useful for high scores in games.
		ChatInstance string `json:"chat_instance"`
		// Optional. Data associated with the callback button. Be aware that a
		// bad client can send arbitrary data in this field.
		Data string `json:"data"`
		// Optional. Short name of a Game to be returned, serves as the unique
		// identifier for the game.
		GameShortName string `json:"game_short_name"`
	}

	// KeyboardButton represents one button of the reply keyboard. For simple
	// text buttons String can be used instead of this object to specify text of
	// the button.
	KeyboardButton struct {
		// Text of the button. If none of the optional fields are used, it will
		// be sent to the bot as a message when the button is pressed
		Text string `json:"text"`
	}

	// ReplyKeyboardMarkup represents a custom keyboard with reply options.
	ReplyKeyboardMarkup struct {
		// Array of button rows, each represented by an Array of KeyboardButton
		// objects.
		Keyboard [][]*KeyboardButton `json:"keyboard"`
		// Optional. Requests clients to resize the keyboard vertically for
		// optimal fit (e.g., make the keyboard smaller if there are just two
		// rows of buttons). Defaults to false, in which case the custom
		// keyboard is always of the same height as the app's standard keyboard.
		ResizeKeyboard bool `json:"resize_keyboard"`
		// Optional. Requests clients to hide the keyboard as soon as it's been
		// used. The keyboard will still be available, but clients will
		// automatically display the usual letter-keyboard in the chat – the
		// user can press a special button in the input field to see the custom
		// keyboard again. Defaults to false.
		OneTimeKeyboard bool `json:"one_time_keyboard"`
		// Optional. Use this parameter if you want to show the keyboard to
		// specific users only. Targets: 1) users that are @mentioned in the
		// text of the Message object; 2) if the bot's message is a reply (has
		// reply_to_message_id), sender of the original message.
		Selective bool `json:"selective"`
	}
)

// Request and response wrappers are defined here.
type (
	Response struct {
		// The response contains a JSON object, which always has a Boolean field
		// ‘ok’ and may have an optional String field ‘description’ with a
		// human-readable description of the result. If ‘ok’ equals true, the
		// request was successful and the result of the query can be found in
		// the ‘result’ field. In case of an unsuccessful request, ‘ok’ equals
		// false and the error is explained in the ‘description’.
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}

	GetMeResponse struct {
		Response
		Result *User `json:"result"`
	}

	GetUpdatesResponse struct {
		Response
		Result []Update `json:"result"`
	}

	SendMessageRequest struct {
		// Unique identifier for the target chat or username of the target
		// channel (in the format @channelusername).
		ChatID int64 `json:"chat_id"`
		// Text of the message to be sent.
		Text string `json:"text"`
		// Send Markdown or HTML, if you want Telegram apps to show bold,
		// italic, fixed-width text or inline URLs in your bot's message.
		ParseMode string `json:"parse_mode"`
		// Disables link previews for links in this message.
		DisableWebPagePreview bool `json:"disable_web_page_preview"`
		// Sends the message silently.
		// iOS users will not receive a notification, Android users will
		// receive a notification with no sound.
		DisableNotification bool `json:"disable_notification"`
		// If the message is a reply, ID of the original message.
		ReplyToMessageID int `json:"reply_to_message_id"`
		// Additional interface options. A JSON-serialized object for an inline
		// keyboard, custom reply keyboard, instructions to hide reply keyboard
		// or to force a reply from the user.
		// TODO: Support InlineKeyboardMarkup, ReplyKeyboardHide and ForceReply.
		ReplyMarkup *ReplyKeyboardMarkup `json:"reply_markup,omitempty"`
	}

	SendMessageResponse struct {
		Response
		Result *Message `json:"result"`
	}

	ForwardMessageRequest struct {
		// Unique identifier for the target chat or username of the target
		// channel (in the format @channelusername).
		ChatID int64 `json:"chat_id"`
		// Unique identifier for the chat where the original message was sent
		// (or channel username in the format @channelusername).
		FromChatID int64 `json:"from_chat_id"`
		// Sends the message silently.
		// iOS users will not receive a notification, Android users will receive
		// a notification with no sound.
		DisableNotification bool `json:"disable_notification"`
		// Unique message identifier.
		MessageID int `json:"message_id"`
	}

	ForwardMessageResponse struct {
		Response
		Result *Message `json:"result"`
	}

	SendStickerRequest struct {
		// Unique identifier for the target chat or username of the target
		// channel (in the format @channelusername).
		ChatID int64 `json:"chat_id"`
		// Sticker to send. Pass a file_id as String to send a file that exists
		// on the Telegram servers (recommended), pass an HTTP URL as a String
		// for Telegram to get a .webp file from the Internet, or upload a new
		// one using multipart/form-data..
		Sticker string `json:"sticker"`
		// Sends the message silently.
		// iOS users will not receive a notification, Android users will receive
		// a notification with no sound.
		DisableNotification bool `json:"disable_notification"`
		// If the message is a reply, ID of the original message.
		ReplyToMessageID int `json:"reply_to_message_id"`
		// Additional interface options. A JSON-serialized object for an inline
		// keyboard, custom reply keyboard, instructions to hide reply keyboard
		// or to force a reply from the user.
		// TODO: ReplyMarkup
	}

	SendStickerResponse struct {
		Response
		Result *Message `json:"result"`
	}
)

// Call Telegram API method.
func (e *Bot) CallMethod(method string, params interface{}) ([]byte, error) {
	url := "https://api.telegram.org/bot" + e.token + "/" + method
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// A simple method for testing your bot's auth token. Requires no parameters.
// Returns basic information about the bot in form of a User object.
func (e *Bot) GetMe() (*User, error) {
	res, err := e.CallMethod("getMe", nil)
	if err != nil {
		return nil, err
	}
	me := &GetMeResponse{}
	err = json.Unmarshal(res, me)
	if err != nil {
		return nil, err
	}
	if !me.OK {
		return nil, errors.New(me.Description)
	}
	return me.Result, nil
}

// Receive incoming updates using long polling (wiki). An Array of Update
// objects is returned.
func (e *Bot) GetUpdates(offset, limit, timeout int) ([]Update, error) {
	res, err := e.CallMethod("getUpdates", map[string]int{
		"offset":  offset,
		"limit":   limit,
		"timeout": timeout,
	})

	if err != nil {
		return nil, err
	}
	updates := &GetUpdatesResponse{}
	err = json.Unmarshal(res, updates)
	if err != nil {
		return nil, err
	}
	if !updates.OK {
		return nil, errors.New(updates.Description)
	}
	return updates.Result, nil
}

// Send text messages. On success, the sent Message is returned.
func (e *Bot) SendMessage(body *SendMessageRequest) (*Message, error) {
	res, err := e.CallMethod("sendMessage", body)
	if err != nil {
		return nil, err
	}
	message := &SendMessageResponse{}
	err = json.Unmarshal(res, message)
	if err != nil {
		return nil, err
	}
	if !message.OK {
		return nil, errors.New(message.Description)
	}
	return message.Result, nil
}

// Forward messages of any kind. On success, the sent Message is returned.
func (e *Bot) ForwardMessage(body *ForwardMessageRequest) (*Message, error) {
	res, err := e.CallMethod("forwardMessage", body)
	if err != nil {
		return nil, err
	}
	message := &ForwardMessageResponse{}
	err = json.Unmarshal(res, message)
	if err != nil {
		return nil, err
	}
	if !message.OK {
		return nil, errors.New(message.Description)
	}
	return message.Result, nil
}

// TODO: sendPhoto

// TODO: sendAudio

// TODO: sendDocument

// Send .webp stickers. On success, the sent Message is returned.
func (e *Bot) SendSticker(body *SendStickerRequest) (*Message, error) {
	res, err := e.CallMethod("sendSticker", body)
	if err != nil {
		return nil, err
	}
	message := &SendStickerResponse{}
	err = json.Unmarshal(res, message)
	if err != nil {
		return nil, err
	}
	if !message.OK {
		return nil, errors.New(message.Description)
	}
	return message.Result, nil
}

// TODO: sendVideo

// TODO: sendVoice

// TODO: sendLocation

// TODO: sendVenue

// TODO: sendContact

// TODO: sendChatAction

// TODO: getUserProfilePhotos

// TODO: getFile

// TODO: kickChatMember

// TODO: leaveChat

// TODO: unbanChatMember

// TODO: getChat

// TODO: getChatAdministrators

// TODO: getChatMembersCount

// TODO: getChatMember

// TODO: answerCallbackQuery

// TODO: editMessageText

// TODO: editMessageCaption

// TODO: editMessageReplyMarkup

// TODO: answerInlineQuery

// TODO: sendGame

// TODO: setGameScore

// TODO: getGameHighScores
