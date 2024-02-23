

function appendFolder(folderName) {
    let inp = document.getElementById("query-folder");
    let v = inp.value + "/" + folderName;
    let fl = document.getElementById("download-folder");
    fl.value = v;
}