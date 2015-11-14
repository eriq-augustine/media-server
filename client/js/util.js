"use strict";

var mediaserver = mediaserver || {};
mediaserver.util = mediaserver.util || {};

mediaserver.util.hashPass = function(pass, username) {
  var salted = mediaserver.util.saltPass(pass, username);
  var hash = CryptoJS.SHA3(salted, {outputLength: 512}).toString(CryptoJS.enc.Hex);
  return hash;
}

mediaserver.util.saltPass = function(pass, username) {
  return username + "." + pass + "." + username;
}

mediaserver.util.addTokenParam = function(link) {
   if (!mediaserver.apiToken) {
      return link;
   }

   var params = {
      "token": mediaserver.apiToken
   };
   return link + '?' + $.param(params);
}
