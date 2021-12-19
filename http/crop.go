package http

//func cropImageHandler(w http.ResponseWriter, r *http.Request, session *interfaces.Session) {
//	img := "." + r.FormValue("img")
//	top, _ := strconv.ParseFloat(r.FormValue("top"), 32)
//	left, _ := strconv.ParseFloat(r.FormValue("left"), 32)
//	bottom, _ := strconv.ParseFloat(r.FormValue("bottom"), 32)
//	right, _ := strconv.ParseFloat(r.FormValue("right"), 32)
//	fType := strings.Split(img, ".")
//
//	// check if the user has edit permission for the model
//	modelNameParts := strings.Split(img, "/")
//	if len(modelNameParts) < 2 {
//		utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "The image path is too short"})
//		return
//	}
//	modelName := modelNameParts[len(modelNameParts)-2]
//	modelName = strings.Split(modelName, "_")[0]
//	//if !session.User.GetAccess(modelName).Edit {
//	//	utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "You don't have permission to edit this model"})
//	//	return
//	//}
//
//	// Open the file
//	imageFile, err := os.Open(img)
//	if err != nil {
//		utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "Unable to open image. Check logs for details"})
//		return
//	}
//	defer imageFile.Close()
//	var myImage image.Image
//	if strings.ToLower(fType[len(fType)-1]) == "jpg" || strings.ToLower(fType[len(fType)-1]) == "jpeg" {
//		myImage, err = jpeg.Decode(imageFile)
//		if err != nil {
//			utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "Unable to decode JPEG image. Check logs for details"})
//			return
//		}
//	}
//	if strings.ToLower(fType[len(fType)-1]) == "png" {
//		myImage, err = png.Decode(imageFile)
//		if err != nil {
//			utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "Unable to decode PNG image. Check logs for details"})
//			return
//		}
//	}
//
//	if strings.ToLower(fType[len(fType)-1]) == "gif" {
//		myImage, err = gif.Decode(imageFile)
//		if err != nil {
//			utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "Unable to decode GIF image. Check logs for details"})
//			return
//		}
//	}
//
//	rect := image.Rect(int(left), int(top), int(right), int(bottom))
//
//	mySubImage := image.NewRGBA(rect)
//	draw.Draw(mySubImage, rect, myImage, image.Point{int(top), int(left)}, draw.Src)
//
//	f, err := os.Create(strings.Replace(img, "_raw", "", -1))
//	if err != nil {
//		utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "Unable to create image. Check logs for details"})
//		return
//	}
//	defer f.Close()
//
//	if strings.ToLower(fType[len(fType)-1]) == preloaded.CJPG || strings.ToLower(fType[len(fType)-1]) == preloaded.CJPEG {
//		err = jpeg.Encode(f, mySubImage, nil)
//		if err != nil {
//			utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "Unable to encode JPEG image. Check logs for details"})
//			return
//		}
//	}
//	if strings.ToLower(fType[len(fType)-1]) == preloaded.CPNG {
//		err = png.Encode(f, mySubImage)
//		if err != nil {
//			utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "Unable to encode PNG image. Check logs for details"})
//			return
//		}
//	}
//
//	if strings.ToLower(fType[len(fType)-1]) == preloaded.CGIF {
//		o := gif.Options{}
//		err = gif.Encode(f, mySubImage, &o)
//		if err != nil {
//			utils.ReturnJSON(w, r, map[string]string{"status": "error", "err_msg": "Unable to encode GIF image. Check logs for details"})
//			return
//		}
//	}
//	utils.ReturnJSON(w, r, map[string]string{"status": "ok"})
//}
//
