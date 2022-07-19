---
title: Hello World!
publish_date: 2022-07-13
---

This is my eight-billionth first blog post! *Yay!!!*

By this point, you could safely call my incessant need to make new blogs an addiction. I've always loved writing. It was my first passion, my first career choice, and my first heartbreak. I've just never found a good blog framework that made writing easy and customization fun. See, the perfect blog infrastructure sits somewhere between **fast-to-deploy** and **easy-to-customize**.

![#](/hello_world_1.jpeg)


Fast-to-deploy blogs are typically managed by someone other than the author. They're not cookie cutter--you can typically change fonts, background color, or post layout. But typically you're unable to really customize how your pieces look. Think Substack or Medium. They usually come with the added bonus of an ecommerce or membership suite, letting you sell access to your writing or advertise on your posts with little to no effort. Typically, they advertise themselves as an all-inclusive package. The writer picks it up, gives it a spit-shine, and never touches it again. And why would they? Most writers don't care about how their blog works. They'd be just as happy writing for a paper, or magazine. That's okay! I'm not here to shame them. But the mindset isn't really right for me.

On the other hand, you've got the very technical writers. Those brave souls who want to customize the whole process. They'll spend hours setting up their WordPress blog, integrating the right plugins, or--if they're really <s>crazy</s> dedicated they can build their whole business [from scratch](https://stratechery.com/).

I'm somewhere in the middle. I don't really enjoy using the fast-to-deploy fully managed platforms. They're fine for handling the monotony of implementing a payment system, or designing a homepage, but when it comes to telling a story, or showing off specifically formatted text--common for developers they fall short. One of my biggest disappointment’s with Substack is how they render code. I'm a web developer who dabbles in embedded programming (think the programs that actually run on your computer). When I want to convey something I've done, I can wrap the code in a formatting block which will strip it of all the flair you and I need for normal text, but isn't very useful for developers.

```js
console.log("Hello, World!")
```

The block adds color to the code, making it easier to read for developers. This formatting implicitly conveys a lot of extra useful information, making my life as a writer and your life as a reader easier. Substack doesn't do this. It strips the plain text formatting, but with no colors and no easy way to manually add them, it's obvious programming authors are an afterthought. 

<!-- add image of substack code -->

"But that's just one market," I hear you say. "Substack caters to intellectuals, and writers who [felt persecuted by cancel culture](https://www.theguardian.com/commentisfree/2020/nov/17/substack-media-platform-publishing)." Damn, I didn't know you were so educated. But you're missing the point. It's not just programming authors that need special formatting. Data-backed journalists, industry reporters, and anyone writing about anything more complicated than their opinion all have the need to display and share complicated information: charts, reports, and demonstrations. Imagine if [Ciechanowski](https://ciechanow.ski/mechanical-watch/) was limited to a content management system (CMS)? The same problem persists on Medium and Svbtle. Any platform that doesn't allow you to get into its code is one you're destined to outgrow. 

## Somewhere in the Middle

Somewhere in the middle is a solution that is easy-to-customize and fast-to-deploy. It requires technical skills, but it's not so uncultivated as to make maintaining it a full-time job. This blog is running on [Deno Deploy](https://deno.com/deploy). A brand-new hosting system full of buzzwords and promises. I'm not totally convinced it's anything new, but it does have a generous free tier atypical of its competitors--if I spend another cent on a new project that goes nowhere, I'll cut off my phalanges. I'm using the [Blog](https://deno.land/x/blog@0.4.1) template to make running this site as painless as possible. On another branch of its Git repository I'm developing a custom solution, but having the flexibility of always being able to switch back to this implementation in under a minute makes it more of a fun hobby instead of "*life or death of the product*".

As far as all of this relates to you: enjoy yourself. My personal goals writing here are that while it is still fun, this will be the final resting place of all my ambiguous digital writing. I won't write here about my local town, or specifically about one of my products unless it can be more broadly applied to you or some idea. What I will write about are the general things that interest me, have propelled me through obscurity to my present stand-still, and bring me real joy. Free from the watchful eye of those we now, immortality live in our mind's eye. If any of my bullshit has you enamored, consider [adding this site](https://bottle.quest/feed) to your RSS reader.

Your's truely,

u8