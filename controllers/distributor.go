package controllers

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	// "math/rand"
	// "time"
	// "encoding/json"
	"sort"
)

type distributorPredictor func(*Distributor) bool

const (
	color_orange = "orange"
	color_red    = "red"
	color_grey   = "grey"
	color_purple = "purple"
)

var error_no_websocket_connection = errors.New("error_no_websocket_connection")

type CheckPoint int

func getCheckPointMap() map[string]CheckPoint {
	m := make(map[string]CheckPoint)
	m["checkpoint_flag_origin"] = checkpoint_flag_origin
	m["checkpoint_flag_prepared_for_game"] = checkpoint_flag_prepared_for_game
	m["checkpoint_flag_game_started"] = checkpoint_flag_game_started
	m["checkpoint_flag_game_over"] = checkpoint_flag_game_over
	// m["checkpoint_flag_order_distribute_over"] = checkpoint_flag_order_distribute_over
	return m
}

var (
	checkpoint_flag_origin            CheckPoint = 0 //只是刚登录系统
	checkpoint_flag_prepared_for_game CheckPoint = 1 //进入游戏房间等待
	checkpoint_flag_game_started      CheckPoint = 2 //游戏已经开始
	checkpoint_flag_game_over         CheckPoint = 3 //游戏结束
	// checkpoint_flag_order_distribute_over CheckPoint = 3
	// checkpoint_max                        CheckPoint = 4
)

// 配送员
type Distributor struct {
	// ID                string
	// Name              string
	// Color             string //地图上marker颜色
	UserInfo          *User
	AcceptedOrders    OrderList
	CheckPoint        CheckPoint //所处的关卡
	Online            bool
	StartPos, DestPos *Position //配送时设置的出发和目的路径点
	CurrentPos        *Position //配送时实时所在的路径
	NormalSpeed       float64   //运行速度 km/h
	CurrentSpeed      float64   //当前运行速度，0表示停止
	Distance          float64   //所在或者将要行驶的路径长度
	line              *Line
	GameTimeMaxLength int //游戏最大时长
	TimeElapse        int //运行时间
	Score             int //得分
	Rank              int //排名
	GameID            string
	Conn              *websocket.Conn `json:"-"` // Only for WebSocket users; otherwise nil.
	// AtGame            bool
	// MaxAcceptedOrdersCount int             `json:"-"` //配送员可以接收的最大订单数量
}

func NewDistributor(user *User) *Distributor {
	return &Distributor{
		UserInfo:       user,
		AcceptedOrders: OrderList{},
		CheckPoint:     checkpoint_flag_origin,
		// ID:             id,
		// Name:           name,
		// Color:          color,
		// MaxAcceptedOrdersCount: maxCount,
	}
}
func (this *Distributor) String() string {
	return fmt.Sprintf("ID: %-3s  Name: %-4s 游戏进程：%d  接收的订单：%2d online:%t score: %d timeElapse: %d rank: %d checkPoint: %d",
		this.UserInfo.ID, this.UserInfo.Name, this.CheckPoint, len(this.AcceptedOrders), this.IsOnline(), this.Score, this.TimeElapse, this.Rank, this.CheckPoint)
}
func (d *Distributor) PosString() string {
	if d.CurrentPos == nil {
		return fmt.Sprintf("ID: %-3s  Name: %-4s  未设定当前位置", d.UserInfo.ID, d.UserInfo.Name)
	}
	if d.DestPos == nil {
		return fmt.Sprintf("ID: %-3s  Name: %-4s  (%f, %f) %fkm/h", d.UserInfo.ID, d.UserInfo.Name, d.CurrentPos.Lng, d.CurrentPos.Lat, d.CurrentSpeed)
	}
	return fmt.Sprintf("ID: %-3s  Name: %-4s  (%f, %f) => (%f, %f) %fkm/h", d.UserInfo.ID, d.UserInfo.Name, d.CurrentPos.Lng, d.CurrentPos.Lat, d.DestPos.Lng, d.DestPos.Lat, d.CurrentSpeed)
}

// //接收了订单数量的上限
// func (this *Distributor) fullyLoaded() bool {
// 	return len(this.AcceptedOrders) >= this.MaxAcceptedOrdersCount
// }

func (d *Distributor) whetherHasEndTheGame() bool {
	return d.CheckPoint == checkpoint_flag_game_over
}

func (d *Distributor) caculateScore() {
	unSignedOrders := d.AcceptedOrders.Filter(func(o *Order) bool { return o.Signed == false })
	d.Score -= unSignedOrders.totalScore(0) * 2
}

//接收配送订单
func (this *Distributor) acceptOrder(order *Order) {
	this.AcceptedOrders = append(this.AcceptedOrders, order)
}

func (d *Distributor) IsOriginal() bool {
	return d.CheckPoint == checkpoint_flag_origin
}
func (d *Distributor) IsOnline() bool {
	return d.Conn != nil
}
func (d *Distributor) GetID() string {
	return d.UserInfo.ID
}
func (d *Distributor) SetConn(conn *websocket.Conn) {
	d.Conn = conn
	d.Online = true
}
func (d *Distributor) SendBinaryMessage(msg []byte) error {
	if d.Conn != nil {
		// DebugTraceF("distributor [%s] send message to client", d.ID)
		return d.Conn.WriteMessage(websocket.TextMessage, msg)
	}
	return error_no_websocket_connection
}
func (d *Distributor) IdEqals(id string) bool {
	return d.UserInfo.ID == id
	// return d.Is(func(dr *Distributor) bool { return dr.ID == id })
}
func (d *Distributor) SetOffline() error {
	defer func() {
		d.Online = false
	}()
	if d.Conn != nil {
		if err := d.Conn.Close(); err == nil {
			d.Conn = nil
			DebugSysF("[%s] OffLine, close websocket ", d.UserInfo.Name)
			// DebugInfoF("[%s] OffLine WebSocket closed", d.ID)
		} else {
			DebugMustF("[%s] OffLine,But close websocket err: %s", d.UserInfo.Name, err)
			return err
		}
	}
	// d.preparedForOrderSelect = false
	return nil
}
func (d *Distributor) setCheckPoint(check CheckPoint) {
	d.CheckPoint = check
}
func (d *Distributor) copyAll() *Distributor {
	return &Distributor{
		UserInfo:       d.UserInfo.copy(),
		AcceptedOrders: OrderList{},
		CheckPoint:     d.CheckPoint,
		// ID:             d.ID,
		// Name:           d.Name,
		// Color: d.Color,
		// MaxAcceptedOrdersCount: d.MaxAcceptedOrdersCount,
	}
}

func (d *Distributor) Is(f distributorPredictor) bool {
	return f(d)
}

type DistributorList []*Distributor

func (dl DistributorList) Rank() DistributorList {
	sort.Sort(dl)
	f := func(d *Distributor, lastRank, lastScore, lastTimeLength int) (myRank, myScore, myTimeLength int) {
		if d.Score == lastScore && d.TimeElapse == lastTimeLength {
			d.Rank = lastRank
		} else {
			d.Rank = lastRank + 1
		}
		return d.Rank, d.Score, d.TimeElapse
	}
	var myRank, myScore, myTimeLength int = 0, -1, -1
	for _, d := range dl {
		// d.Rank = i + 1
		myRank, myScore, myTimeLength = f(d, myRank, myScore, myTimeLength)
	}
	return dl
}
func (dl DistributorList) Len() int {
	return len(dl)
}
func (dl DistributorList) Swap(i, j int) {
	dl[i], dl[j] = dl[j], dl[i]
}
func (dl DistributorList) Less(i, j int) bool {
	if dl[i].Score == dl[j].Score {
		return dl[i].TimeElapse < dl[j].TimeElapse
	} else {
		return dl[i].Score > dl[j].Score

	}
}
func (dl DistributorList) clone(f distributorPredictor) (l DistributorList) {
	for _, d := range dl {
		if f == nil || f(d) {
			l = append(l, d.copyAll())
		}
	}
	return
}
func (dl DistributorList) forEach(f func(*Distributor)) {
	if f == nil {
		return
	}
	for _, d := range dl {
		f(d)
	}
}
func (dl DistributorList) forOne(f distributorPredictor) {
	if f == nil {
		return
	}
	for _, d := range dl {
		if f(d) {
			return
		}
	}
}

func (dl DistributorList) every(f distributorPredictor) bool {
	for _, d := range dl {
		if f(d) == false {
			return false
		}
	}
	return true
}
func (dl DistributorList) filter(f distributorPredictor) (l DistributorList) {
	for _, d := range dl {
		if f == nil || f(d) {
			l = append(l, d)
		}
	}
	return
}
func (dl DistributorList) findOne(f distributorPredictor) *Distributor {
	for _, p := range dl {
		if f(p) {
			return p
		}
	}
	return nil
}
func (this DistributorList) ListName() string {
	return "配送员信息"
}
func (this DistributorList) InfoList() (list []string) {
	for _, d := range this {
		list = append(list, d.String())
	}
	return
}
