<script setup lang="ts">
import {reactive} from 'vue';
import {Buffer} from "buffer";
import {Upload} from '../rest.js';
import {getAllFiles, addFile, clearFiles} from '../indexedDB'

let files = []

const data = reactive({
  uploaded: {
    status: 0,
    id : '',
    name : '',
    size : 0,
    path : '',
    type : '',
    data : null,
  },
  response: {}
});

function uploadSucceeded(res) {
  data.response = res
}

function uploadFailed(error) {
  data.response = error
}

function uploadFile(event) {

  //erase placeholder
  document.getElementById('placeholder').innerHTML = ''
  let upload = new Upload()
  let file = event.target.files[0]
  let temp_uploaded = {
    status: 0,
    id : '',
    name : '',
    size : 0,
    path : '',
    type : '',
    data : null,
  }

  upload.addFile(file)

  upload.onsuccess = uploadSucceeded
  upload.onfail = uploadFailed

  upload.submit()

  temp_uploaded.type = file.type
  
  if(temp_uploaded.type != 'image/png' && temp_uploaded.type != 'image/jpeg'){
    alert('Only PNG and JPG files are allowed')
    return
  }else{
    const today = new Date();
    temp_uploaded.id = today.getFullYear()+'-'+(today.getMonth()+1)+'-'+today.getDate()+'+'+file.name;
    temp_uploaded.name = file.name
    temp_uploaded.size = file.size
    var buf = new Promise(getBuffer(file))
    buf.then(function(data){
      temp_uploaded.data = data
      console.log(temp_uploaded)
      files.push(temp_uploaded)
      addFile(temp_uploaded)
      displayFiles()
      document.getElementById('placeholder').innerHTML = `<img style="max-width: 100%; max-height: 25rem; object-fit: contain;" class="frame" src="data:${temp_uploaded.type};base64,${temp_uploaded.data}" alt=${temp_uploaded.name} />`
    }).catch(function(error){
      console.log("Error: ",error)
    })
  }
}

function getBuffer(file) {
  return function(resolve){
    var reader = new FileReader();
    reader.readAsArrayBuffer(file);
    reader.onload = function() {
      var arrayBuffer = reader.result;
      var bytes: any
      if(typeof arrayBuffer === 'string') {
        bytes = new Uint8Array(arrayBuffer.length);
      } else {
        var raw = new Uint8Array(arrayBuffer);
        // var decoeder = new TextDecoder('utf-8');
        bytes = encodeBase64(raw);
      }
      resolve(bytes);
    }
  }
}

function encodeBase64(data) {
  return Buffer.from(data).toString('base64');
}

function displayFiles() {
  const list = document.getElementById('list')
  if(files != null){
    files.forEach(element => {
      // <img width="25%" src={`data:${IFPipe.type};base64,${IFPipe.data}`} alt={IFPipe.name} />
      list.innerHTML += '<li><div class="frame"><div class="thumbnail"><img class="thumbnail" src="data:'+element.type+';base64,'+element.data+'" alt='+element.name+'></div><div class="text"><h1>Test</h1><p>testing</p></div><div></li>'
    });
  }
}

window.onload = () => {
  getAllFiles().then(function(result:any){
    console.log(result)
    files = result
    displayFiles()
  })
}

function clearCache() {
  files = []
  document.getElementById('placeholder').innerHTML = ''
  document.getElementById('list').innerHTML = ''
  clearFiles()
}

</script>

<template>
  <div id="container">
    <h2>Upload Images</h2>
    <div>
      <label class="custom-file-upload button">
        <input id="upload" accept="image/png,image/jpeg" class="button" type="file" @change="uploadFile">
        <i class="fa fa-cloud-upload"></i> Select File
      </label>
      {{ data.response }}
      
      <div id="placeholder"></div>
      <span class="previous" >Previous Uploads</span><a href="" class="clear_btn" @click="clearCache">clear</a>
      <ul id="list"></ul>
    </div>
  </div>
</template>

<style scoped>

input[type="file"] {
  display: none;
}

</style>

<style>
.thumbnail {
  max-width: 10vw;
  max-height: 10vh;
  object-fit: contain;
  flex-basis: 40%;
}
.frame {
  background-color: rgb(44, 44, 46);
  border-radius: 5px;
  padding: 5px;
  margin: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.text {
  padding-left: 20px;
}
#container {
  justify-content: center;
}
/* keep left justified */
#placeholder {
  height: 25rem;
}
/* keep right justified */
.clear_btn {
  float: right;
  color: rgb(108, 194, 2);
  text-decoration: none;
}
.previous {
  font-size: 24px;
  font-weight: bold;
}
</style>
