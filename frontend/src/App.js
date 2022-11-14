import TopBar from './components/TopBar';
import Uploader from './components/Uploader';
import Storage from './components/Storage';
import { useEffect, useState } from 'react';
import { Container } from 'react-bootstrap';
import axios from 'axios';
function App() {
  const [hashes, setHashes] = useState([])

  useEffect(() => {
    axios.get("http://localhost:8080/l").then(e => {
      setHashes(e.data)
    })
  }, [])

  function handleUpload() {
    console.log('uploaded file...');
    axios.get("http://localhost:8080/l").then(e => {
      setHashes(e.data)
    })
  }

  return (

    

    <div>
      <TopBar />
      <Container className="pt-4">
        <Uploader onUpload={handleUpload} />
        <hr/>
        <Storage hashes={hashes} />
      </Container>
    </div>
  );
}

export default App;
