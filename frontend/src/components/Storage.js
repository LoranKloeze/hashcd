import Card from 'react-bootstrap/Card';
import ListGroup from 'react-bootstrap/ListGroup';
import axios from 'axios';
import { useState } from 'react';

import { Button } from 'react-bootstrap';

function humanFileSize(size) {
  var i = size === 0 ? 0 : Math.floor(Math.log(size) / Math.log(1024));
  return (size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
}

function Storage({ hashes = [] }) {
  const [headers, setHeaders] = useState([])

  function retrieveFile(hash) {
    axios.get(`http://localhost:8080/d/${hash}`).then(e => {
      setHeaders(e.headers)
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
          <pre>{headerItems}</pre>

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
