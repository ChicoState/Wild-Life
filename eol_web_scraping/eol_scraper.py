#EOL is the encyclopedia of Life
#This script scrapes 116 images of western pacific oak from their website

from bs4 import BeautifulSoup 
try: import requests
except ImportError:
    print("Error importing requests")
    exit
url = 'https://eol.org/pages/582277/media?page='
file = 'eol.html'
#requests is for getting the html page from the web
for i in list(range(1, 6)): 
    page_url = url + str(i)
    r = requests.get(page_url)
    if(r.status_code == 200): 
#BeautifulSoup is for parsing the html page
        soup = BeautifulSoup(r.text, "html.parser")
        images = soup.find_all("img", alt="Image of western poison-oak")
        for image in images: 
            img_url = image.find("src", alt="Image of western poison-oak")
            img_url = image.get('src')
            print(img_url.strip())
    else:
        print(page_url + "does NOT exist")
        break

