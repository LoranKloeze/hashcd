import Card from 'react-bootstrap/Card';
import ListGroup from 'react-bootstrap/ListGroup';
import HashInfo from './HashInfo';

function Storage({ hashes = [] }) {

  const listItems = hashes.map((hash) => <HashInfo hash={hash}/>)

  return (
    <div>
      <Card>
        <Card.Body>
          <Card.Title>Hashes in storage</Card.Title>
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
