import Card from 'react-bootstrap/Card';
import Spinner from 'react-bootstrap/Spinner';
import ListGroup from 'react-bootstrap/ListGroup';
import axios from 'axios';
import { useState } from 'react';

import { Button } from 'react-bootstrap';


const devEnv = !process.env.NODE_ENV || process.env.NODE_ENV === 'development'

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
  const [activeHash, setActiveHash] = useState(null)

  function infoFile(hash) {
    setLoading(true)
    axios.get(`/d/${hash}?t=${new Date().getTime()}`).then(e => {
      setHeaders(e.headers)
      setLoading(false)
      setActiveHash(hash)
    })
  }

  const headerItems = Object.keys(headers).map((k) =><div key={k}>{k}: {headers[k]}</div>)

  const listItems = hashes.map((hash) =>
    <ListGroup.Item key={hash.hash}>
      <div class="row align-items-start">
        <div class="col-2">
          <Button variant='light' className="me-2" size='sm' onClick={() => infoFile(hash.hash)}>
            Info
          </Button>
          <a href={devEnv ? `http://localhost:8080/d/${hash.hash}` : `/d/${hash.hash}`} target="_blank" rel="noreferrer" className="btn btn-sm btn-primary">Download</a>
        </div>
        <div class="col-10">
          <pre className="mb-0">{hash.hash} | {humanFileSize(hash.size)}</pre>
          {hash.hash === activeHash && <pre className="mt-2">{headerItems}</pre>}
        </div>
      </div>

      
    </ListGroup.Item>
  )

  return (
    <div>
      <Card>
        <Card.Body>
          <Card.Title>Hashes in storage</Card.Title>
          {loading && <Waiting />}
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
