Bot `@MurmanNewsBot` scraps content of **Murmansk region's** news:      
* https://hibiny.com  
* https://hibinform.ru  
* https://severpost.ru  

and weather: 
* https://pogoda51.ru/
     
and post it to `@murmnews` telegram channel   
_User_ can suggest news to _bot_, suggested news are sent to _admin_, _admin_ approve it and news are posted with notification to _user_ who suggested news.  

Place `config.yaml` file in application directory with content:  

    tgToken: your_token  
    tgChannel: '@channel_name'_ 
    tgAdmin: admin_id  
    tgChatID: channel_chat_id 
    tgSocks: socks5_ip:socks5_port
    tgSocksUser: socks_user 
    tgSocksPass: socks_pass 

`tgChannel` is news channel  
`tgAdmin` is chatID of channel's admin   
`tgChatID` is chatID of channel  

Flags:   
`-bg`  **run in background**   
`-socks`  **run with socks5**   

Core libs used in project:
 * **Scraper** `https://github.com/gocolly/colly`
 * **Telegram API** `https://github.com/go-telegram-bot-api/telegram-bot-api`
 * **Key-Value Storage** `https://github.com/tidwall/buntdb`