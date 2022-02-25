import cv2 as cv 
import sys 

filename = "tommy.jpg"

img = cv.imread(cv.samples.findFile(filename))

if img is None: 
    sys.exit("Could not read the image")

cv.imshow("Display window", img)
k = cv.waitKey(0)

if k == ord("s"):
    cv.imwrite(filename, img)
