function uploadFile() {
    var fileInput = document.getElementById('fileInput');
    var file = fileInput.files[0];

    if(file) {
        console.log("File upload on website click on Upload")
    }

    if (!file) {
        alert('Please select a file.');
        return;
    }

    var formData = new FormData();
    formData.append('file', file);

    fetch('http://localhost:3000/csv/updateFromCSV', {
        method: 'PUT',
        body: formData
    })
    .then(data => {
        document.getElementById('result').innerText = "File Uploaded to database";
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('result').innerText = 'Error uploading file.';
    });
}


//Login

function submitForm() {
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;

    const formData = new FormData();
    formData.append("username", username);
    formData.append("accountnum", password);
}
