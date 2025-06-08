use reqwest::multipart;
use std::error::Error;
use std::fs::File;
use std::io::{self, Write};
use std::path::PathBuf;
use std::str::FromStr;
use zip::ZipWriter;
use zip_extensions::write::ZipWriterExtensions;

#[tokio::main]
async fn main() {
    /* Zip Folder and send it
    zip_guest_code("guest-code/.".to_string()).await.unwrap();
    send_archive("coole-image-id".to_string()).await.unwrap();
    */
    /* request guest code by image id and unzip it  */
    get_archive("coole-image-id".to_string()).await.unwrap();
    unwrap_archive("archive.zip".to_string(), "test-folder".to_string())
        .await
        .unwrap();
}
// param: path_to_folder = which folder to zip
async fn zip_guest_code(path_to_folder: String) -> io::Result<()> {
    let path = PathBuf::from_str(&path_to_folder).unwrap();
    let file = File::create("archive.zip")?;

    // zip_dir returns its second parameter wrapped in a Result if successful
    let zip = ZipWriter::new(file);
    zip.create_from_directory(&path)?;

    Ok(())
}
// archive_path = absolute or relative path from root directory to archive / zip
// target_dir = path to wanted folder for extraction
pub async fn unwrap_archive(
    archive_path: String,
    target_dir: String,
) -> Result<(), Box<dyn Error>> {
    let archive_file = PathBuf::from_str(&archive_path).unwrap();
    let target = PathBuf::from_str(&target_dir).unwrap();
    zip_extensions::zip_extract(&archive_file, &target)?;

    Ok(())
}

pub async fn send_archive(image_id: String) -> Result<(), Box<dyn Error>> {
    let url = format!("http://localhost:6005/proofs?image_id={}", image_id);

    let form = multipart::Form::new().file("file", "archive.zip").await?; // this is now required

    let client = reqwest::Client::new();
    let response = client.post(&url).multipart(form).send().await?;

    println!("Status: {}", response.status());
    println!("Body: {}", response.text().await?);

    Ok(())
}

pub async fn get_archive(image_id: String) -> Result<(), Box<dyn Error>> {
    let url = format!("http://localhost:6005/proofs/{}", image_id);

    let client = reqwest::Client::new();
    let response = client.get(&url).send().await?;

    println!("Status: {}", response.status());

    let bytes = response.bytes().await?;

    let mut file = File::create("archive.zip").unwrap();
    file.write_all(&bytes).unwrap();

    println!("Archive saved as archive.zip");

    Ok(())
}
