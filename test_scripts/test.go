package main

import (
	"fmt"
	"math/rand"
    "time"
	"os"
	"github.com/eatmoreapple/openwechat"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var r *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
}

func generateRandomString(length int) string {
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

func GetRandInt64(n int64) int64 {
	return r.Int63n(n)
}

func KeepAlive(bot *openwechat.Self) {

	ticker := time.NewTicker(time.Minute * 120)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			heartBeat(bot)
		}
	}
}

func heartBeat(bot *openwechat.Self) {
	// 向文件传输助手发送消息，不要再关注公众号了
	// 生成要发送的消息
	outMessage := fmt.Sprintf("防微信自动退出登录[%d]", GetRandInt64(1500))
	bot.SendTextToFriend(openwechat.NewFriendHelper(bot), outMessage)
}

func createFolderIfNotExists(folderPath string) error {
    _, err := os.Stat(folderPath)
    if os.IsNotExist(err) {
        return os.MkdirAll(folderPath, os.ModePerm)
    }
    return nil
}

func main() {

	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	// bot := openwechat.DefaultBot()
	defer reloadStorage.Close()
	err := bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())

	// bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// save_folder is the date of time 
	default_name :=  "default" 
	save_folder :=  "test_image"
	// print the save_folder 
	fmt.Println(save_folder) 
	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
		//  get the user name
		friend_user , err := msg.Sender()
		if err != nil { 
			fmt.Println(err)
			return
		}
		
		default_name = friend_user.NickName 

		if msg.IsPicture(){
			// rand.Seed(time.Now().UnixNano())
			randomImageName := generateRandomString(20)
			// 拼接图像后缀，这里以.jpg 为例
			imageNameWithSuffix := randomImageName + ".jpg"
			file_path := save_folder +"/" + time.Now().Format("2006-01-02") +  "/" + default_name
			// create the folder if not exists
			fmt.Println("file_path:", file_path) 
			err := createFolderIfNotExists(file_path)
                if err!= nil {
                    fmt.Println("创建文件夹失败：", err)
                    return
                }
			msg.SaveFileToLocal( file_path+"/" + imageNameWithSuffix)
		}
			

	}
	// 注册登陆二维码回调
	// bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// // 登陆
	// if  err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// // 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)
	go KeepAlive(self)
	// KeepAlive(self)
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
