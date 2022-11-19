import Card from 'react-bootstrap/Card';
import Form from 'react-bootstrap/Form';
import ListGroup from 'react-bootstrap/ListGroup';
import axios from 'axios';
import { useState } from 'react';

function Uploader({onUpload = () => null}) {
  const [finished, setFinished] = useState(false)
  const [hash, setHash] = useState('')
  const [status, setStatus] = useState('')

  function uploadFile(e) {
    e.preventDefault();
    console.log('uploading file...');

    const formData = new FormData();

    formData.append(
      "f",
      e.target.files[0],
      e.target.files[0].name
    );

    axios.post("http://localhost:8080/u", formData).then(e => {
      onUpload()
      setFinished(true)
      setHash(e.data.hash)
      console.log(e)

      switch (e.status) {
        case 201:
          setStatus('File created')
          break;
        case 200:
          setStatus('File was already uploaded')
          break;
        default:
          setStatus('Something went wrong: ' + e.statusText)
          break;
      }
    })

  }

  return (
    <div>
      <Card>
        <Card.Body>
          <Card.Title>Upload file</Card.Title>
          <Form.Group className="mb-3" controlId="formBasicEmail">
            <Form.Control type="file" onChange={uploadFile} placeholder="Enter email" />
          </Form.Group>
        </Card.Body>
        {finished &&
          <ListGroup className="list-group-flush">
            <ListGroup.Item>File uploaded!</ListGroup.Item>
            <ListGroup.Item className="small bg-success text-light">{status}</ListGroup.Item>
            <ListGroup.Item><pre>{hash}</pre></ListGroup.Item>
          </ListGroup>
        }
      </Card>
    </div>
  );
}

export default Uploader;
