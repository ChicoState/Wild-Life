url_file = "eol_img_urls" 
urls = [] 
with open(url_file, 'r') as f: 
    urls = f.readlines()
unique_urls = set(urls)
for unique_url in unique_urls: 
    print(unique_url)
print(len(unique_urls))
