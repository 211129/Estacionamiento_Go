package scenes

import (
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
	"image/color"
	"math/rand"
	"parking-concurrency/models"
	"sync"
	"time"
)

var (
	spots = []*models.Place{
		models.NewPlace(380, 70, 410, 100, 1, 1),
		models.NewPlace(425, 70, 455, 100, 1, 2),
		models.NewPlace(470, 70, 500, 100, 1, 3),
		models.NewPlace(515, 70, 545, 100, 1, 4),
		models.NewPlace(560, 70, 590, 100, 1, 5),
		models.NewPlace(380, 160, 410, 190, 2, 6),
		models.NewPlace(425, 160, 455, 190, 2, 7),
		models.NewPlace(470, 160, 500, 190, 2, 8),
		models.NewPlace(515, 160, 545, 190, 2, 9),
		models.NewPlace(560, 160, 590, 190, 2, 10),
		models.NewPlace(380, 250, 410, 280, 3, 11),
		models.NewPlace(425, 250, 455, 280, 3, 12),
		models.NewPlace(470, 250, 500, 280, 3, 13),
		models.NewPlace(515, 250, 545, 280, 3, 14),
		models.NewPlace(560, 250, 590, 280, 3, 15),
		models.NewPlace(380, 340, 410, 370, 4, 16),
		models.NewPlace(425, 340, 455, 370, 4, 17),
		models.NewPlace(470, 340, 500, 370, 4, 18),
		models.NewPlace(515, 340, 545, 370, 4, 19),
		models.NewPlace(560, 340, 590, 370, 4, 20),
	}
	parking    = models.NewParkingLot(spots)
	doorMutex  sync.Mutex
	carManager = models.NewOrchestrator()
)

type ParkingLotScene struct {
}

func NewParkingLotScene() *ParkingLotScene {
	return &ParkingLotScene{}
}

func (ps *ParkingLotScene) Start() {
	isFirstTime := true

	_ = oak.AddScene("parkingScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			setUpScene(ctx)

			event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
				if !isFirstTime {
					return 0
				}

				isFirstTime = false

				for {
					go run(ctx)

					time.Sleep(time.Millisecond * time.Duration(getRandomNumber(1000, 2000)))
				}

				return 0
			})
		},
	})
}

func setUpScene(ctx *scene.Context) {

	parkingArea := floatgeom.NewRect2(0, 0, 1000, 1000)
	backgroundImg, _ := render.LoadSprite("assets/images/background.jpeg")
	entities.New(ctx, entities.WithRect(parkingArea), entities.WithRenderable(backgroundImg), entities.WithDrawLayers([]int{0}))

	borderColor := color.RGBA{128, 0, 0, 255}
	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(340, 5, 630, 10)), entities.WithColor(borderColor), entities.WithDrawLayers([]int{0}))
	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(340, 400, 630, 405)), entities.WithColor(borderColor), entities.WithDrawLayers([]int{0}))
	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(340, 70, 345, 400)), entities.WithColor(borderColor), entities.WithDrawLayers([]int{0}))
	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(625, 10, 630, 400)), entities.WithColor(borderColor), entities.WithDrawLayers([]int{0}))

	for _, spot := range spots {
		entities.New(ctx, entities.WithRect(floatgeom.NewRect2(spot.GetArea().Min.X(), spot.GetArea().Min.Y(), spot.GetArea().Min.X()+2.5, spot.GetArea().Max.Y())), entities.WithColor(color.RGBA{212, 172, 13, 255}))
		entities.New(ctx, entities.WithRect(floatgeom.NewRect2(spot.GetArea().Max.X(), spot.GetArea().Min.Y(), spot.GetArea().Max.X()-2.5, spot.GetArea().Max.Y())), entities.WithColor(color.RGBA{212, 172, 13, 255}))
	}
}

func run(ctx *scene.Context) {
	car := models.NewCar(ctx)

	carManager.Add(car)

	car.Enqueue(carManager)

	spotAvailable := parking.GetPlaceAvailable()

	doorMutex.Lock()

	car.JoinDoor(carManager)

	doorMutex.Unlock()

	car.Park(spotAvailable, carManager)

	time.Sleep(time.Millisecond * time.Duration(getRandomNumber(40000, 50000)))

	car.LeavePlace(carManager)

	parking.SetFreePlace(spotAvailable)

	car.Leave(spotAvailable, carManager)

	doorMutex.Lock()

	car.ExitDoor(carManager)

	doorMutex.Unlock()

	car.GoLong(carManager)

	car.Remove()

	carManager.Remove(car)
}

func getRandomNumber(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}
