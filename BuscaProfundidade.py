import random
import copy
import networkx as nx
from matplotlib import pyplot as plt

def gerarTabuleiro(v=10):
  slots = [x for x in map(lambda a: a if a!=0 else None, range(9))]
  #random.shuffle(slots)
  #return [ slots[x*3:x*3+3] for x in range(3) ]
  base = [ slots[x*3:x*3+3] for x in range(3) ]
  for x in range(v):
      child = generateChildStates(base)
      base = child[random.randint(0,len(child)-1)]
  return base

def isDone(tabuleiro):
  return tabuleiro == [[None, 1 ,2],[3,4,5],[6,7,8]]

def stateHash(tabuleiro):
  coef = [2**x for x in range(1,10)]
  result = 0
  idx = 0
  for x in tabuleiro:
    for y in x:
      result += coef[idx]* (y if y else 0)
      idx += 1
  return result

def generateChildStates(tabuleiro):
  X,Y = 0,0
  for idy, y in enumerate(tabuleiro):
    for idx, x in enumerate(y):
      if x == None:
        X = idx
        Y = idy
  new_states = []
  for x in [X-1,X+1]:
    if x < 0 or x > 2:
      continue
    new_states.append(copy.deepcopy(tabuleiro))
    new_states[-1][Y][X] = new_states[-1][Y][x]
    new_states[-1][Y][x] = None
  for y in [Y-1,Y+1]:
    if y < 0 or y > 2:
      continue
    new_states.append(copy.deepcopy(tabuleiro))
    new_states[-1][Y][X] = new_states[-1][y][X]
    new_states[-1][y][X] = None
  return new_states

def buscaEmProfundidade(tabuleiro, G):

    states_visited = []
    state_stack = [copy.deepcopy(tabuleiro)]
    state_stack_count = [-1]
    while len(state_stack) > 0:
      current = state_stack[-1]
      
      if not stateHash(current) in states_visited:
        state_stack_count[len(state_stack)-1] += 1
        states_visited.append(stateHash(current))
        G.add_node(states_visited.index(stateHash(current)), pos=(state_stack_count[len(state_stack)-1], len(state_stack)))
        if len(states_visited) >= 2:
            G.add_edge(states_visited.index(stateHash(state_stack[-2])), states_visited.index(stateHash(current)) )
        
      
      if isDone(current):
        return (True, current)

      state_finished = True
      for x in generateChildStates(current):
        if not stateHash(x) in states_visited:
          state_stack_count.append(state_stack_count[len(state_stack)-1]-1)
          state_stack.append(x)
          state_finished = False
          break
      if state_finished:
        state_stack.pop()
    return False

G = nx.Graph()
print(buscaEmProfundidade(gerarTabuleiro(4), G))

pos=nx.get_node_attributes(G,'pos')
nx.draw(G,pos,font_weight='bold')
#plt.savefig("Teste.pdf", dpi=2500)
plt.show()