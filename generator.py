# Import necessary functions and properties from libraries
from py_ecc.bn128 import is_on_curve, FQ  # Functions to check if a point is on the curve and to handle field elements
from py_ecc.fields import field_properties  # To get properties of the field used in bn128
field_mod = field_properties["bn128"]["field_modulus"]  # Get the modulus of the field for bn128
from hashlib import sha256  # Import SHA-256 hashing function
from libnum import has_sqrtmod_prime_power, sqrtmod_prime_power  # Functions to check and find square roots in modular arithmetic

b = 3  # Coefficient for the elliptic curve equation y^2 = x^3 + b
#For bn128 b=3

#Bulletproofs rely on an existing, standardized elliptic curve (like BN128, BLS12-381, or secp256k1).
seed = "BulletProofs"  # Seed value for generating a starting point
#publicly agreed-upon string 

# Generate a starting x value by hashing the seed and reducing it modulo the field modulus
convert_string_to_bytes=seed.encode('ascii')
sha_256_bytes=sha256(convert_string_to_bytes)
hexadecimal_conversion=sha_256_bytes.hexdigest()
integer_conversion=int(hexadecimal_conversion,16)
#x=integer_conversion%field_mod
x = int(sha256(seed.encode('ascii')).hexdigest(), 16) % field_mod 

entropy = 0  # Initialize entropy counter
vector_basis = []  # List to store generated points on the curve
# Loop to find a valid point on the elliptic curve
while not has_sqrtmod_prime_power((x**3 + b) % field_mod, field_mod, 1) :
    # Increment x to find a point that is on the curve
    x = (x + 1) % field_mod  # Wrap around if x exceeds field modulus
    entropy = entropy + 1  # Increase entropy count
    
# Determine y value based on the parity of entropy (even or odd) 
#sqrt generates 2 values of y so that not same sign is considered every time
#sign is decided on the basis of entropy
y = list(sqrtmod_prime_power((x**3 + b) % field_mod, field_mod, 1))[entropy & 1 == 0]
# Create a point (x, y) in the field
point = (FQ(x), FQ(y))  # Convert x and y to field elements
# Check if the point is on the curve; raise an error if not
assert is_on_curve(point, b), "sanity check"  # Ensure the point is valid

# Add the valid point to the vector basis
vector_basis.append(point) #We are collecting G1,G2....

# Generate a new x value by hashing the current x and reducing it modulo the field modulus
x = int(sha256(str(x).encode('ascii')).hexdigest(), 16) % field_mod 
# Print the list of points generated
print(vector_basis)  # Output the points found on the curve

