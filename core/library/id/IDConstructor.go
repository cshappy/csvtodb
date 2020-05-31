package id

import "time"

var (
	machineID     int64 //机器的ID
	sn            int64 //随机数
	lastTimeStamp int64 // 时间戳
)

/**
machineId 当前机器的编号
 */
func CreateId(machine int64) int64{
	machineID = machine << 12
	// 获取上一个ID生成的时间
	lastTimeStamp = time.Now().UnixNano() / 1000000

	// 获取当前ID生成的时间
	var currentTimeStamp = time.Now().UnixNano() / 1000000
	// 用一时间多人创建
	if currentTimeStamp == lastTimeStamp {
		sn ++
		// 预留的12个bit位不足
		if sn > 4095 {
			time.Sleep(time.Millisecond)
			currentTimeStamp = time.Now().UnixNano() / 1000000
			lastTimeStamp = currentTimeStamp
			sn = 1
		}
	}

	if currentTimeStamp > lastTimeStamp {
		sn = 1
		lastTimeStamp = currentTimeStamp
	}
	if currentTimeStamp < lastTimeStamp {
		return 0
	}

	rightBinValue := currentTimeStamp & 0x1FFFFFFFFFF
	rightBinValue <<= 22
	return rightBinValue | machineID |sn
}
