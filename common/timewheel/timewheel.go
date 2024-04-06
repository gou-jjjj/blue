package timewheel

import (
	"container/list"
	"time"
)

// location 结构体用于存储任务在时间轮中的位置信息
type location struct {
	slot  int
	etask *list.Element
}

// TimeWheel 是一个能够定时执行任务的时间轮结构体
type TimeWheel struct {
	interval time.Duration // 时间轮的转动周期
	ticker   *time.Ticker  // 定时器，用于触发时间轮的转动
	slots    []*list.List  // 时间轮的槽，每个槽是一个链表，存储等待执行的任务

	timer             map[string]*location // 任务映射表，用于快速定位任务
	currentPos        int                  // 当前时间轮的位置
	slotNum           int                  // 时间轮的槽数量
	addTaskChannel    chan task            // 添加任务的通道
	removeTaskChannel chan string          // 删除任务的通道
	stopChannel       chan bool            // 停止时间轮的通道
}

// task 结构体定义了一个任务的属性
type task struct {
	delay  time.Duration // 任务的延迟时间
	circle int           // 任务需要在时间轮上转的圈数
	key    string        // 任务的唯一标识
	job    func()        // 任务执行的函数
}

// New 创建一个新的时间轮实例
// interval: 时间轮的转动周期
// slotNum: 时间轮的槽数量
func New(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}

	tw := &TimeWheel{
		interval:          interval,
		slots:             make([]*list.List, slotNum),
		timer:             make(map[string]*location),
		currentPos:        0,
		slotNum:           slotNum,
		addTaskChannel:    make(chan task),
		removeTaskChannel: make(chan string),
		stopChannel:       make(chan bool),
	}
	tw.initSlots()

	return tw
}

// initSlots 初始化时间轮的槽，每个槽初始化为一个空的链表
func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

// Start 启动时间轮，开始转动
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

// Stop 停止时间轮的转动
func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

// AddJob 添加一个新任务到时间轮上
// delay: 任务的延迟时间
// key: 任务的唯一标识
// job: 任务执行的函数
func (tw *TimeWheel) AddJob(delay time.Duration, key string, job func()) {
	if delay < 0 {
		return
	}

	tw.addTaskChannel <- task{delay: delay, key: key, job: job}
}

// RemoveJob 从时间轮上移除一个任务
// key: 任务的唯一标识
func (tw *TimeWheel) RemoveJob(key string) {
	if key == "" {
		return
	}
	tw.removeTaskChannel <- key
}

// start 是一个goroutine，负责处理时间轮的转动、任务的添加和移除
func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(&task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

// tickHandler 处理时间轮每次转动的逻辑
func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
	go tw.scanAndRunTask(l)
}

// scanAndRunTask 遍历当前槽的任务列表，执行到期的任务，并从列表中移除
func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					panic(err)
				}
			}()
			job := task.job
			job()
		}()
		next := e.Next()
		l.Remove(e)
		if task.key != "" {
			delete(tw.timer, task.key)
		}
		e = next
	}
}

// addTask 将一个任务添加到时间轮上
// task: 需要添加的任务
func (tw *TimeWheel) addTask(task *task) {
	pos, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle

	e := tw.slots[pos].PushBack(task)
	loc := &location{
		slot:  pos,
		etask: e,
	}
	if task.key != "" {
		_, ok := tw.timer[task.key]
		if ok {
			tw.removeTask(task.key)
		}
	}
	tw.timer[task.key] = loc
}

// getPositionAndCircle 根据任务的延迟时间计算任务应该放置的槽的位置和需要在时间轮上转的圈数
// 返回值: pos - 任务所在槽的位置；circle - 任务需要在时间轮上转的圈数
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle = delaySeconds / intervalSeconds / tw.slotNum
	pos = (tw.currentPos + delaySeconds/intervalSeconds) % tw.slotNum

	return
}

// removeTask 根据任务的唯一标识从时间轮上移除任务
// key: 任务的唯一标识
func (tw *TimeWheel) removeTask(key string) {
	pos, ok := tw.timer[key]
	if !ok {
		return
	}
	l := tw.slots[pos.slot]
	l.Remove(pos.etask)
	delete(tw.timer, key)
}
