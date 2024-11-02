package mqtt

// func TestPublishSite(t *testing.T) {
// 	pkgconfig.SetConfigFile("../../config.ini")
// 	bootstrap.Register(
// 		pkgconfig.Init,
// 		config.Init,
// 		zap.Init,
// 		logrus.Init,
// 		metrics.Init,
// 		Init,
// 	)
//
// 	assert.NoError(t, bootstrap.Init())
//
// 	count := 0
// 	for {
// 		r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 		Publish("/site/site_01", model.MessageSiteCommand{
// 			SiteId:         "site_" + util.UUID(),
// 			SiteType:       "pad",
// 			CommandUuid:    util.UUID(),
// 			CommandTime:    strconv.Itoa(int(time.Now().Unix())),
// 			CommandAction:  model.ActionIn,
// 			CommandCargoId: "cargo_" + util.UUID(),
// 			SiteCode:       strconv.Itoa(r.Intn(100)),
// 			InsiteCode:     strconv.Itoa(r.Intn(100)),
// 			OutsiteCode:    strconv.Itoa(r.Intn(100)),
// 		})
// 		count++
// 		if count > 10 {
// 			break
// 		}
// 		time.Sleep(time.Second)
// 	}
// }
//
// func TestPublishAgv(t *testing.T) {
// 	pkgconfig.SetConfigFile("../../config.ini")
// 	bootstrap.Register(
// 		pkgconfig.Init,
// 		config.Init,
// 		zap.Init,
// 		logrus.Init,
// 		metrics.Init,
// 		Init,
// 	)
//
// 	assert.NoError(t, bootstrap.Init())
// 	var wg sync.WaitGroup
// 	publishCount := 10
// 	total := 10000
// 	wg.Add(publishCount)
// 	for i := 1; i <= publishCount; i++ {
// 		go func() {
// 			agvId := fmt.Sprintf("agv_%03d", i)
// 			defer wg.Done()
// 			count := 0
// 			for {
// 				r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 				x := r.Intn(100)
// 				y := r.Intn(100)
// 				yaw := r.Intn(100)
// 				Publish(fmt.Sprintf("/agv/%s/status", agvId), model.MessageAgvStatus{
// 					AgvId:              agvId,
// 					AgvType:            "agv_notype",
// 					CurrentTime:        strconv.Itoa(int(time.Now().Unix())),
// 					CurrentLocation:    fmt.Sprintf("%v", []int{x, y, yaw}),
// 					CurrentLocationX:   strconv.Itoa(x),
// 					CurrentlocationY:   strconv.Itoa(y),
// 					CurrentLocationYaw: strconv.Itoa(yaw),
// 					CurrentCargoId:     "cargo_" + util.UUID(),
// 					CurrentLiftStatus:  model.RandomLiftStatus(),
// 					CurrentStatus:      model.RandomAgvStatus(),
// 					TaskSite:           "site_" + util.UUID(),
// 					TaskUuid:           util.UUID(),
// 					TaskAction:         model.RandomAction(),
// 					RunArea:            strconv.Itoa(r.Intn(20)),
// 				})
// 				count++
// 				if count > total {
// 					break
// 				}
// 				time.Sleep(time.Second)
// 			}
// 		}()
// 	}
//
// 	wg.Wait()
// }
