
# Import necessary functions and properties from libraries
from py_ecc.bn128 import is_on_curve, FQ  # Functions to check if a point is on the curve and to handle field elements
from py_ecc.fields import field_properties  # To get properties of the field used in bn128
field_mod = field_properties["bn128"]["field_modulus"]  # Get the modulus of the field for bn128
from hashlib import sha256  # Import SHA-256 hashing function
from libnum import has_sqrtmod_prime_power, sqrtmod_prime_power  # Functions to check and find square roots in modular arithmetic
from py_ecc.bn128 import G1, multiply, add, FQ
import random
from py_ecc.bn128 import curve_order as p

def generate_points_on_curve(b,seed,n):
    # Generate a starting x value by hashing the seed and reducing it modulo the field modulus
    x = int(sha256(seed.encode('ascii')).hexdigest(), 16) % field_mod 
    entropy=0
    vector_basis=[]
    while len(vector_basis)<n:
        while not (has_sqrtmod_prime_power((x**3+b)%field_mod,field_mod,1)):
            x=(x+1)%field_mod
            entropy+=1
        y=list(sqrtmod_prime_power((x**3+b)%field_mod,field_mod,1))[entropy&1==0] 
        point=(FQ(x),FQ(y))
        assert is_on_curve(point, b), "sanity check"  # Ensure the point is valid
        vector_basis.append(point) #We are collecting G1,G2....
        x = int(sha256(str(x).encode('ascii')).hexdigest(), 16) % field_mod #generate new x for next iteration
    return vector_basis

def random_field_element():
    return random.randint(0, p)

def commit(G,B,f,gamma) :
    C=add(multiply(G,f),multiply(B,gamma))
    return C



b = 3  # Coefficient for the elliptic curve equation y^2 = x^3 + b
seed = "BulletProofs"  # Seed value for generating a starting point
#publicly agreed-upon string 
n=4
inally=generate_points_on_curve(b,seed,2*n+2)
G_vect=[]
H_vect=[]
G_vect.append(inally[0])
G_vect.append(inally[1])
G_vect.append(inally[2])
G_vect.append(inally[3])
H_vect.append(inally[4])
H_vect.append(inally[5])
H_vect.append(inally[6])
H_vect.append(inally[7])
G=inally[8]
B=inally[9]
# print("G is ",G ,"\n B is",B,"\n Vector G is ",G_vect,"\n Vector H is ",H_vect);
alpha=random_field_element()
beta=random_field_element()
gamma=random_field_element()
gamma1=random_field_element()
gamma2=random_field_element()
A=commit()