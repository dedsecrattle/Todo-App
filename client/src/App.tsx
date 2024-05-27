import useSWR from 'swr'
import { Box ,  List,
  ListItem,
  Heading,
  Text,
  Icon,
  Flex,
  Center,
  Button,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,} from '@chakra-ui/react';
import { MdCheckBox, MdRadioButtonChecked, MdRadioButtonUnchecked } from 'react-icons/md';
import { useState } from 'react';

export const ENDPOINT = "http://localhost:4000";

interface Todo {
  id : number,
  title : string,
  body : string,
  done : boolean
}

const fetcher = (url: string) =>
  fetch(`${ENDPOINT}/${url}`).then((r) => r.json());

function App() {
  const {data, mutate} = useSWR<Todo[]>("todo", fetcher)

  const [isOpen, setIsOpen] = useState(false);
  const [formData, setFormData] = useState({ title: "", body: ""});

  const handleOpen = () => setIsOpen(true);
  const handleClose = () => setIsOpen(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const requestOption = {
      method : "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body : JSON.stringify({
        "title": formData.title,
        "body": formData.body,
        "done": false
      })
    }
    const response = await fetch(`${ENDPOINT}/todo`, requestOption).then((r) => r.json())
    const updatedData : Todo[] = [...data || [], {id: response.id ,title: response.title, body: response.body, done: response.done}]
    mutate(updatedData, false)
    handleClose();
  };

  const handleDoneClick = async (id: number) => {
    const updatedData = await fetch(`${ENDPOINT}/todo/${id}`, {method: "PATCH"}).then((r) => r.json())
    mutate(updatedData, false)
  }

  return (
    <Flex minH="100vh" bg="gray.100" align="center" justify="center">
      <Center display='flex' flexDirection='column'>
        <Box bg="white" p={6} rounded="md" boxShadow="md" minW="xl" maxW='xl'>
          <Heading as="h1" size="xl" textAlign="center" mb={4}>
            Todos
          </Heading>
          <List>
          {data?.map((todo) => {
          return <ListItem display='flex' key={`todo_list__${todo.id}`}>
          <Center gap='5' padding='5'>
           {todo.done ? 
           <Icon boxSize='8' as={MdRadioButtonChecked} color='green.500' onClick = {() => handleDoneClick(todo.id)}></Icon>
          :<Icon boxSize='8' as={MdRadioButtonUnchecked} color='red.500' onClick = {() => handleDoneClick(todo.id)}></Icon>} 
           <Text fontSize='30px'>{todo.title}</Text> 
           <Text fontSize='20px' border='1px' padding='5'>{todo.body}</Text>
           </Center>
          </ListItem>})}
          </List>
        </Box>
        <Button bg='red.400' padding='25px' color='white' margin='20px' onClick={handleOpen}> Add Todo</Button>
    <Modal isOpen={isOpen} onClose={handleClose}>
    <ModalOverlay />
    <ModalContent>
      <ModalHeader>Add New Todo</ModalHeader>
      <ModalCloseButton />
      <ModalBody>
        <form onSubmit={handleSubmit}>
          <Flex direction="column" mb={4}>
            <Text mb={2}>Title:</Text>
            <Input
              name="title"
              value={formData.title}
              onChange={handleChange}
            />
          </Flex>
          <Flex direction="column" mb={4}>
            <Text mb={2}>Body:</Text>
            <Input
              name="body"
              value={formData.body}
              onChange={handleChange}
            />
          </Flex>
          <ModalFooter>
            <Button colorScheme="blue" mr={3} type="submit">
              Submit
            </Button>
            <Button variant="ghost" onClick={handleClose}>
              Cancel
            </Button>
          </ModalFooter>
        </form>
      </ModalBody>
    </ModalContent>
  </Modal>
      </Center>
    </Flex>
  )
}

export default App
