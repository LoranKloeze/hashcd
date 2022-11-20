import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';

function TopBar() {
  return (
    <div>
      <Navbar bg="light" expand="lg">
      <Container>
        <Navbar.Brand href="#home">HashCD</Navbar.Brand>
        
      </Container>
    </Navbar>
    </div>
  );
}

export default TopBar;
