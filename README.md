Telegram Bots that might help developers to work as extra tools
 

## Requirements
* _tasks_rpg_bot_
  * github.com/mattn/go-sqlite3 - db read/write tasks
  * github.com/joho/godotenv - configuring

 
## List
* _echo_bot_ - echos messages sent to bot, formats them, if command is given
* _tasks_rpg_bot_ - mark, review, remind tasks list in fun RPG game-like way
 
 
## TODO
* check permissions for accessing/using bots
* use for telegram bots https://github.com/go-telegram-bot-api/telegram-bot-api or similar
* add forms for RPG tasks bot https://core.telegram.org/bots#keyboards
* listeners for transactions - add/edit tasks in multiple steps etc.
* use go routines and channels for each type of command to process message updates
* messageUpdate -> fan out channels x Commands -> fan out each message entity -> fan in message entities -> fan in command Done status
* instead of crappy solutions for SendMessageStruct and ToArray objects make a proper one with structs and JSON flags