import sys
try: import requests
except ImportError:
    print("Error importing requests")
    exit
try: import shutil
except ImportError:
    print("Failed to import shutil")
    exit
#Open file of urls and read them boi 
url_file = "eol_img_urls" 
urls = [] 
with open(url_file, 'r') as f: 
    urls = f.readlines()
img_name = "western_poison_oak"
img_extension = '.jpg'
iteration = 0
file_out = ""

#Fetch the image and write it to disk 
for url in urls: 
    res = requests.get(url.strip(), stream=True)
    if(res.status_code == 200): 
        file_out = img_name + str(iteration) + img_extension
        with open(file_out, 'wb') as out_file: 
            res.raw.decode_content = True
            shutil.copyfileobj(res.raw, out_file)
        iteration += 1
    else:
        print("Failed to get", str(url))
    del res

def get_one_image(): 
    url='https://content.eol.org/data/media/7f/ff/3d/542.4547261920.130x130.jpg'
    url='https://content.eol.org/data/media/7f/ff/3d/542.4547261920.130x130.jpg'
    res = requests.get(url, stream=True)
    if(res.status_code == 200): 
        print(res.raw)
        with open('yeehaw.jpg', 'wb') as out_file: 
            res.raw.decode_content=True
            shutil.copyfileobj(res.raw, out_file)
    else: 
        print("nope")
    del res
    sys.exit()


