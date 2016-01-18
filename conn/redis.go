package conn

import (
	"github.com/garyburd/redigo/redis"
)

//All redis actions

func SetMasterId(id int) {
	c := Pool.Get()
	defer c.Close()
	c.Do("SET", "evolsnowChatId", id)
}

func GetMasterId() int {
	c := Pool.Get()
	defer c.Close()
	id, _ := redis.Int(c.Do("GET", "evolsnowChatId"))
	return id
}

func SetUserChatId(user string, id int) {
	c := Pool.Get()
	defer c.Close()
	key := user + "ChatId"
	c.Do("SET", key, id)
}

func GetUserChatId(user string) int {
	c := Pool.Get()
	defer c.Close()
	key := user + "ChatId"
	id, _ := redis.Int(c.Do("GET", key))
	return id
}

func HSetMemo(time, user, memo string) {
	c := Pool.Get()
	defer c.Close()
	var setMemoLua = `
	local id = redis.call("INCR", "memoIncrId")
	redis.call("RPUSH", user..":memo:", id)
	redis.call("HMSET", "memo:"..id, "time", KEYS[1], "content", KEYS[2])
	`
	script := redis.NewScript(2, setMemoLua)
	script.Do(c, time, memo)
}

func HGetAllMemos(user string) (ret []interface{}) {
	c := Pool.Get()
	defer c.Close()
	var multiGetMemoLua = `
	local data = redis.call("LRANGE", KEYS[1]..":memo")
	local ret = {}
  	for idx=1, #data do
  		ret[idx] = redis.call("HGETALL", "memo:"..data[idx])
  	end
  	return ret
   `
	script := redis.NewScript(1, multiGetMemoLua)
	ret, _ = redis.Values(script.Do(c, user))
	return
}

//
//var multiGetScript = redis.NewScript(0, multiGetMemoLua)
