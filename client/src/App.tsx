import useSWR, { mutate } from 'swr'
import { Box ,  List,
  ListItem,
  ListIcon,
  OrderedList,
  Heading,
  CheckboxIcon,
  Text,
  Icon,} from '@chakra-ui/react';
import { MdCheckBox, MdRadioButtonUnchecked } from 'react-icons/md';

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
  return (
    <Box display= 'flex' flexDirection='column' justifyContent='center' alignItems='center' margin='1.5rem'>
      <Heading marginBottom='1.5rem'>
        Todos
      </Heading>
      <List>
      {data?.map((todo) => {
        return <ListItem border='1px' padding='1rem' borderRadius='10' margin='1rem'>
          <Box display='flex' flexDirection='row' w='75%' gap='10' justifyContent='center' alignItems='center' textAlign='center'>
          {todo.done ? <Icon boxSize={8} as={MdCheckBox} color='green.500'></Icon> 
          : <Icon boxSize={8} as={MdRadioButtonUnchecked} color='red.500'></Icon> 
          }
          <Heading size='xl'>{todo.title}</Heading>
          <Text fontSize='3xl'>{todo.body}</Text>
          </Box>
        </ListItem>
      })}
      </List>
    </Box>
  )
}

export default App
