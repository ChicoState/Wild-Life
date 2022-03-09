import {Buffer} from "buffer"
import {fileType} from "./types"

export function getBuffer(file: File) {
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

function encodeBase64(data: Uint8Array) {
   return Buffer.from(data).toString('base64');
}

export function displayFiles(files: fileType[]) {
  const list = document.getElementById('list')
  if(files != null){
    files.forEach(element => {
      list.innerHTML += '<li><div class="frame"><div class="thumbnail"><img class="thumbnail" src="data:'+element.type+';base64,'+element.data+'" alt='+element.name+'></div><div class="text"><h1>Test</h1><p>testing</p></div><div></li>'
    });
  }
}