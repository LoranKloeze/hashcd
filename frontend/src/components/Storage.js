import Card from 'react-bootstrap/Card';
import Spinner from 'react-bootstrap/Spinner';
import ListGroup from 'react-bootstrap/ListGroup';
import axios from 'axios';
import { useState } from 'react';

import { Button } from 'react-bootstrap';



const Waiting = () => {
  return (
    <div className='text-center'>
      <Spinner animation="border" role="status">
        <span className="visually-hidden">Loading...</span>
      </Spinner>
    </div>

  );

}

function humanFileSize(size) {
  var i = size === 0 ? 0 : Math.floor(Math.log(size) / Math.log(1024));
  return (size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
}

function Storage({ hashes = [] }) {
  const [headers, setHeaders] = useState([])
  const [loading, setLoading] = useState(false)

  function retrieveFile(hash) {
    setLoading(true)
    axios.get(`http://localhost:8080/d/${hash}?t=${new Date().getTime()}`).then(e => {
      setHeaders(e.headers)
      setLoading(false)
    })
  }

  const listItems = hashes.map((hash) =>
    <ListGroup.Item key={hash.hash}>
      <Button variant='light' size='sm' onClick={() => retrieveFile(hash.hash)}>
        <pre className="mb-0">{hash.hash} | {humanFileSize(hash.size)}</pre>
      </Button>
    </ListGroup.Item>
  )

  const headerItems = Object.keys(headers).map((k) => {
    if (k === 'content-length') {
      return <div key={k}>{k}: {humanFileSize(headers[k])}</div>
    } else {
      return <div key={k}>{k}: {headers[k]}</div>
    }

  })

  return (
    <div>
      <Card>
        <Card.Body>
          <Card.Title>Hashes in storage</Card.Title>
          { loading ? <Waiting /> : <pre>{headerItems}</pre> }
        </Card.Body>
        {listItems.length > 0 &&
          <ListGroup className="list-group-flush">
            {listItems}
          </ListGroup>
        }
      </Card>
    </div>
  );
}

export default Storage;
