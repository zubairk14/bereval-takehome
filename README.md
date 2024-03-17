Prompt:

A coffee shop usually has several moving pieces: grinders, brewers, different types of beans, and even the number of baristas themselves (plus how fast they can brew coffee- which itself is a function of how fast the grinders + brewers work). What is being given to you is a basic representation (that doesn't work too well in its current form) of a coffee shop, and we want to fix it up so that it behaves more like a real coffee shop. 

 
Some questions/tips to get you started:
- What models should you modify or add (i.e. maybe a `Barista` model)?
- How can we measure that your logic is working correctly? For example, if we only have 2 baristas, we shouldn't be able to handle 1000 requests to make coffee all at the same time- we probably need some order queueing system and signaling that the order is ready later on.
- Different coffee shops have different number of baristas, grinders, etc.- as well as offering different types and strengths of coffee. If the solution would allow for configuration of these different aspects, that would be awesome.
- The code has quite a few comments on where to fix things / make improvements, feel free to follow those comments first and come back to improving different parts later.
 
There is no "exact" solution to this problem, it is only an attempt at modeling a real-life scenario as best as possible (while also making it flexible to model different types of coffee shops).

