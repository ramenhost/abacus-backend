var queryStr = window.location.search;
var id = queryStr.substr(1).split("=")[1]

//Send update POST request with this id in the url
console.log(id)
