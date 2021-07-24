---
"lang": "en",
"title": "MEng Molecular Bioengineering Final Report",
"subtitle": "PyQMI: a python library for receptive field mapping through mutual information optimisation",
"authors": ["Max Taylor-Davies<sup>1</sup>"],
"adresses": ["<sup>1</sup>Department of Bioengineering, Imperial College London"],
"date": "June 2021",
"description": "",
"tags": ["Information theory","Optimisation","Visual system","Receptive fields"]
---

### Abstract

We present an effective and performant GPU-accelerated python implementation of a mutual information-based receptive field estimator. We assess the ability of the estimator to recover simple receptive fields from model visual neurons, and compare performance against established approaches. We then apply the estimator to neural data from the Allen large-scale mouse visual coding dataset, as an example of how it may be used in neuroinformatics analysis.

<br>
<hr>
<br>

- [1. Introduction](#1-introduction)
  - [A. Receptive fields and their estimation](#a-receptive-fields-and-their-estimation)
  - [B. Information-based approaches to RF mapping](#b-information-based-approaches-to-rf-mapping)
- [2. Methods](#2-methods)
  - [A. Computing quadratic mutual information and its derivatives](#a-computing-quadratic-mutual-information-and-its-derivatives)
    - [QMI](#qmi)
    - [First gradient](#first-gradient)
    - [Second gradient](#second-gradient)
  - [B. Gradient-based optimisation of quadratic mutual information](#b-gradient-based-optimisation-of-quadratic-mutual-information)
  - [C. Implementation](#c-implementation)
    - [Modified objective function](#modified-objective-function)
    - [Separate optimisation loop](#separate-optimisation-loop)
  - [D. Testing against simple model neurons](#d-testing-against-simple-model-neurons)
  - [E. Application to the Allen Visual Coding dataset](#e-application-to-the-allen-visual-coding-dataset)
- [3. Results](#3-results)
  - [A. Optimisation approaches](#a-optimisation-approaches)
    - [Comparing algorithms](#comparing-algorithms)
    - [Incorporating RF smoothness](#incorporating-rf-smoothness)
  - [B. Recovery of model neuron RFs](#b-recovery-of-model-neuron-rfs)
  - [C. Analysis of Allen dataset](#c-analysis-of-allen-dataset)
- [4. Discussion](#4-discussion)
- [References](#references)

<br>
<hr>
<br>

## 1. Introduction

### A. Receptive fields and their estimation

First coined by Sherrington in 1906 [[1](#1)], the term <i>receptive field</i> (RF) was originally used to refer to an area of the body surface where a stimulus could elicit a reflex response. The concept was later extended to different types of sensory neurons, and is now understood to describe, for a given sensory neuron, the region of some multidimensional stimulus space that, when stimulated, produces a response. Knowledge of a particular neuron's RF can yield insight into what functional role the neuron might fulfil, and how it relates to other neurons within the same population and at higher and lower levels of processing. RFs of the visual system were originally taken to delineate a two-dimensional area of visual input space; more recently they have been expanded to account for temporal variation, and as such are often referred to as <i>spatiotemporal receptive fields</i>.

Dominant historical approaches to RF estimation can be grouped under the <i>spike-triggered analysis</i> class of methods. In these methods, it is assumed that the response of a neuron at any given point in time is determined solely by the values of the presented stimulus over some recent temporal window. Specifically, the relationship between stimulus and response is assumed to follow an inhomogeneous Poisson process, whose instantaneous frequency is given by applying a linear filter (or set of filters) to the preceding stimulus window, and passing the output through a nonlinear scalar function. The simplest of these methods is <i>spike-triggered average</i> (STA), which assumes a single filter RF and estimates it as the weighted mean of all spike-eliciting stimuli presented:

$$
\begin{equation}
\text{STA} = \frac{1}{n}\sum_{i=1}^T\mathbf{y}_i\mathbf{x}_i
\end{equation}
$$

where $n$ is the total number of spikes, $\mathbf{x}_i$ is the stimulus vector over the $i$th window, and $\mathbf{y}_i$ is the spike count for the $i$th window. STA has been used to characterise RFs of cells in the retinal ganglion [[2](#2)], laterate geniculate nucleus (LGN) [[3](#3)], and primary visual cortex (V1) [[4](#4)]. <i>Spike-triggered covariance</i> (STC) is a similar method, and can be seen as an extension of STA to the case of multiple filters. Assuming a zero-mean stimulus, the STC matrix is computed by

$$
\begin{equation}
\text{STC} = \frac{1}{n}\sum_{i=1}^T\mathbf{y}_i[\mathbf{x}_i - \text{STA}][\mathbf{x}_i - \text{STA}]^T
\end{equation}
$$

If $\mathbf{C}$ is the covariance matrix of the stimulus itself, then eigenvectors of $\text{STC} - \mathbf{C}$ with significantly positive eigenvalues correspond to excitatory filters, and eigenvectors with significantly negative eigenvalues correspond to inhibitory filters.

Despite the appealing simplicity in form and ease of implementation presented by the techniques of spike-triggered analysis, they suffer from numerous limitations that cannot be ignored. The accuracy of STA and STC estimates have strong dependence on the statistical features of the stimuli used. For STA, the estimator is guaranteed to be consistent if and only if the stimuli obey a spherically symmetrical distribution. For STC, the stimuli must be Gaussian-distributed. For simple experiments with synthetic stimuli, this problem can be sidestepped by tailoring the stimulus distributions to satisfy these constraints. However, this is not possible if we want to use stimuli that are more representative of what the visual system would typically be presented with in its <q>natural environment</q>, which may be necessary to form RF estimates that are biologically meaningful. A second limitation of spike-triggered techniques lies in their relative statistical inefficiency, which, in the absence of large quantities of stimulus/response data, can require averaging over time to produce a purely spatial RF estimate, thus losing any information encoded by the temporal dynamics of the neuron's response function.

### B. Information-based approaches to RF mapping

Mutual information (MI) was originally defined by Shannon, based on his definition of entropy. In simplest terms, the MI between two random variables $X$ and $Y$ offers a way to describe the \textit{decrease in uncertainty} about $X$ gained by observing $Y$. There has been a good deal of research into using information-theoretic measures as objective/cost functions in optimisation and machine learning problems. For instance, MI makes sense as an objective function in feature extraction problems, where there is a question of how to project a high-dimensional input to a much lower-dimensional, possibly nonlinear subspace, while preserving the predictive relationship between inputs and class labels. This can be solved by finding the projection that maximises the mutual information between the projected inputs and the class labels.

The idea of taking an information-theoretic approach to the problem of RF estimation was first introduced in 2004 by Sharpee [[5](#5)]. She characterised the receptive field as a set of vectors $\{w^{(1)}, ..., w^{(K)}\}$ also known as the <i>relevant subspace</i> (RS), such that the probability of a neural response depends only on the projections $y^{(\mu)} = w^{(\mu)} \cdot \mathbf{s}, \ \mu = 1, ..., K,$ of the stimulus $\mathbf{s}$ onto the RS. The RS is then estimated through an optimisation procedure: first, some initial vector $v$ is selected, and the information quantity $I(v)$ is computed, serving as an invariant measure of how much the response is determined by stimulus projection onto the particular direction $\hat{v}$. She then optimises over all possible $v$ to find the vector maximising $I(v)$. If the optimal value of $I(v)$ is less than the total spike information $I_\text{spike}$, then the RS must consist of more than one vector, and so the optimisation is repeated over a set of two vectors $(v_1, v_2)$. The number of directions is incremented until $I(v_1, ..., v_n) \approx I_\text{spike}$, at which point the procedure stops, as the set of vectors $(v_1, ..., v_n)$ sufficiently defines the RS. The method was shown to provide more accurate receptive field estimates than STA. Crucially, it was shown to perform well on even strongly non-Gaussian stimuli, such as natural images.

This optimisation procedure is effectively using as its objective function the MI between the projections $Y = \{y^{(\mu)}\}$ and the binary class labels $C = \{$<q>spiking</q>, <q>non-spiking</q>$\}$, given by

$$
\begin{equation}
I(C, S) = \sum_{c \in C}\int_{y^{(\mu)} \in Y}p(c, y^{(\mu)})\log\frac{p(c, y^{(\mu)})}{P(c)P(y^{(\mu)})}
\end{equation}
$$

The reason that MI has not seen more widespread usage as an objective function, and the disadvantage of Sharpee's approach, lies in its inherent computational difficulties - not only are the probability densities of both random variables required, but the numerical integration of these functions introduces significant computational complexity. However, when the objective is not to compute an accurate value of the mutual information between two variables, but rather to find a distribution that maximises or minimises it, it is not actually necessary to use the original formulation (also known as the Kullback-Leibler divergence). For example, an alternative formulation substituting Shannon entropy with Renyi entropy can be used in combination with Parzen density estimation to obtain a non-parametric measure that no longer requires knowledge of the original probability densities. Torkkola [[6](#6)] takes this idea further, presenting a quadratic approximation to mutual information, QMI, that can be used in place of traditional MI as a computationally efficient non-parametric objective function. The problem of RF estimation is then a matter of solving the optimisation problem

$$
\begin{equation}
\hat{\mathbf{w}} = \argmax_{\mathbf{w}} I_Q(\mathbf{y}, \mathbf{r}), \ \mathbf{y} = \mathbf{w} \cdot \mathbf{s}
\end{equation}
$$

This approach has indeed been explored by a small number of researchers. [[9](#9)] was the first instance, and focused on eight retinal ganglion cell (RGC) types from the mouse visual system. Here again, the estimation was framed as a binary classification problem, with class labels 'spiking' and 'non-spiking'. [[10](#10)] then applied the QMI method to cells in the mouse dLGN, demonstrating superior RF estimation relative to STA for datasets consisting of 300 or more stimulus frames. The QMI method was shown to possess greater statistical efficiency than STA, requiring fewer samples to estimate a receptive field. This improved efficiency enabled the extraction of 51\% receptive fields from a previously-published dataset. However, despite these successes, there does not yet exist a high-quality open-source software toolkit for utilising the power of QMI-based RF estimation methods in an intuitive and accessible way. In this work, we aim to develop and validate an effective, performant GPU-accelerat

## 2. Methods

### A. Computing quadratic mutual information and its derivatives

#### QMI

Let $\mathbf{r}$ be a vector representing, for one cell, average spiking frequency at each stimulus frame. If we normalise this vector ($\mathbf{r} = \mathbf{r} / \max(\mathbf{r})$), then we can think of it as representing the prior probabilities for the positive <q>spiking</q> class: $P(c^+) = \mathbf{r}$. The equivalent prior probabilities for the negative <q>non-spiking</q> class are then given by $P(c^-) = 1 - \mathbf{r}$. We can substitute these prior probabilities into Torkkola's original formulation of the QMI, replacing his generalised priors of $P(c_p) = \frac{J_p}{N}$, to obtain the joint probability distributions

$$
\begin{equation}
p(c^+,y) = \frac{1}{N}\sum_{i=1}^Nr_iG(y - y_i, \sigma^2I) \\ p(c^-, y) = \frac{1}{N}\sum_{i=1}^N(1-r_i)G(y-y_i, \sigma^2I)
\end{equation}
$$

These allow us to derive an expression for the information potential $V_\text{IN}$ as

$$
\begin{equation}
\begin{aligned}
V_\text{IN} = \sum_c\int_yp(c,y)^2dy = \int_yp(c^+,y)^2dy + \int_yp(c^-,y)^2dy \\
= \int_y\bigg[\frac{1}{N}\sum_ir_iG(y-y_i,\sigma^2I)\bigg]^2dy + \int_y\bigg[\frac{1}{N}\sum_i(1-r_i)G(y-y_i, \sigma^2I)\bigg]^2dy \\ = \frac{1}{N^2}\int_y\sum_i\sum_jr_ir_jG(y-y_i,\sigma^2I)G(y-y_j,\sigma^2I)dy \\ + \frac{1}{N^2}\int_y\sum_i\sum_j(1-r_i)(1-r_j)G(y-y_i,\sigma^2I)G(y-y_j,\sigma^2I)dy
\end{aligned}
\end{equation}
$$

Using the fact that $\int_yG(y-a,\Sigma_1)G(y-b,\Sigma_2)dy = G(a-b,\Sigma_1+\Sigma_2)$, we can rewrite this as

$$
\begin{equation}
\begin{aligned}
V_\text{IN} = \frac{1}{N^2}\sum_i\sum_jr_ir_jG(y_i-y_j,2\sigma^2I) + \frac{1}{N^2}\sum_i\sum_j(1-r_i)(1-r_j)G(y_i-y_j,2\sigma^2I) \\
= \frac{1}{N^2}[r_ir_j + (1-r_i)(1-r_j)]G(y_i-y_j,2\sigma^2I)
\end{aligned}
\end{equation}
$$

Or, expressed in matrix form,

$$
\begin{equation}
V_\text{IN} = \frac{1}{N^2}\text{sum}\bigg(\text{sum}\bigg([\mathbf{r}'\mathbf{r} + (1-\mathbf{r})'(1-\mathbf{r})] \cdot G(\mathbf{Y})\bigg)\bigg), \ \mathbf{Y}_{ij} = y_i - y_j
\end{equation}
$$

Carrying out a similar process for the second information potential, $V_\text{ALL}$, we find

$$
\begin{equation}
V_\text{ALL} = \frac{1}{N^4}\text{sum}\bigg(\text{sum}\bigg(\bigg[\sum_i\sum_j[r_ir_j + (1-r_i)(1-r_j)]\bigg] \cdot G(\mathbf{Y})\bigg)\bigg)
\end{equation}
$$

and finally for the third potential we have

$$
\begin{equation}
V_\text{BTW} = \frac{1}{N^3}\text{sum}\bigg(\text{sum}\bigg(\mathbf{A} \cdot G(\mathbf{Y})\bigg)\bigg), \ A_{ij} = \sum_{k=1}^Nr_kr_j + (1-r_k)(1-r_j)
\end{equation}
$$

The value of the QMI is then given by

$$
\begin{equation}
    I_Q = V_\text{IN} + V_\text{ALL} - 2V_\text{BTW}
\end{equation}
$$

#### First gradient

We have expressed the QMI, following Torkkola [[6](#6)], as a linear combination of three information potentials. Each of these information potentials is computed as a nested summation over the elementwise product of some matrix and $\mathbf{G} = G(\mathbf{Y})$. Since each of these matrices depends only on the normalised response vector $\mathbf{r}$, and not on the estimated receptive field vector or any other mutable parameter, we can extract them out into a combined weight matrix and rewrite the QMI as

$$
\begin{equation}
I_Q = \text{sum}\big(\text{sum}\big(\mathbf{M} \cdot \mathbf{G} \big )\big), \ \mathbf{M} = \mathbf{M}_\text{IN} + \mathbf{M}_\text{ALL} - 2\mathbf{M}_\text{BTW}
\end{equation}
$$

We can then compute the gradient vector of the QMI with respect to our receptive field vector estimate $\mathbf{w}$ as

$$
\frac{\partial I_Q}{\partial \mathbf{w}} = \frac{\partial}{\partial \mathbf{w}}\bigg[\sum_i\sum_j \mathbf{M}_{ij} \cdot \mathbf{G}_{ij} \bigg] \\ \ \\
= \sum_i\sum_j\bigg[\mathbf{M}_{ij}\cdot\bigg(\frac{\partial \mathbf{G}}{\partial \mathbf{w}}\bigg)_{ij}\bigg]
$$

$\frac{\partial \mathbf{G}}{\partial \mathbf{w}}$ is obtained by

$$
\frac{\partial \mathbf{G}}{\partial \mathbf{w}} = \frac{\partial G}{\partial \mathbf{Y}}\frac{\partial\mathbf{Y}}{\partial\mathbf{w}} \\ \ \\
= \bigg[-\frac{1}{2\sigma^2}\mathbf{Y}\cdot \mathbf{G}\bigg]\cdot \mathbf{S}
$$

where $\mathbf{S}$ is a difference tensor for the extracted stimulus window $\mathbf{s}$, i.e.

$$\mathbf{S} \in \mathbb{R}^{n\times n\times FH}, \ \mathbf{S}_{ij} = \mathbf{s}_i - \mathbf{s}_j$$

and $s_i$ is a vector containing the concatenated flattened stimulus for $H$ frames prior to the response. Weighting by $\mathbf{M}$ and summing over the first two dimensions of the resulting tensor, we have finally

$$
\begin{equation}
\frac{\partial I_Q}{\partial \mathbf{w}} = \sum_i\sum_j\Big(\mathbf{M}_{ij} \cdot \bigg[-\frac{1}{2\sigma^2}\mathbf{Y}_{ij}\cdot \mathbf{G}_{ij}\bigg]\cdot \mathbf{S}_{ij}\Big)
\end{equation}
$$

#### Second gradient

Performing a second differentiation, we can obtain the Hessian matrix $\frac{\partial^2\text{QMI}}{\partial\mathbf{w}^2} \in \mathbb{R}^{FH \times FH}$. The $(a,b)$th element of this matrix can be determined as

$$
\frac{\partial^2I_Q}{\partial\mathbf{w}_a\partial\mathbf{w}_b} = \frac{\partial}{\partial\mathbf{w}_a}\bigg\{\sum_i\sum_j-\frac{\mathbf{M}_{ij}}{2\sigma^2}(\mathbf{Y}_{ij}\mathbf{G}_{ij})\mathbf{S}_{ijb}\bigg\} \\ = \sum_i\sum_j\bigg\{-\frac{\mathbf{M}_{ij}\mathbf{S}_{ijb}}{2\sigma^2}\frac{\partial}{\partial\mathbf{w}_a}(\mathbf{Y}_{ij}\mathbf{G}_{ij})\bigg\} \\
= \sum_i\sum_j\bigg\{-\frac{\mathbf{M}_{ij}\mathbf{S}_{ijb}}{2\sigma^2}\bigg(\mathbf{Y}_{ij}\frac{\partial\mathbf{G}_{ij}}{\partial\mathbf{w}_a} + \mathbf{G}_{ij}\frac{\partial\mathbf{Y}_{ij}}{\partial\mathbf{w}_a}\bigg)\bigg\} \\
= \sum_i\sum_j\bigg\{ -\frac{\mathbf{M}_{ij}\mathbf{G}_{ij}}{2\sigma^2}\bigg(1 - \frac{\mathbf{Y}_{ij}^2}{2\sigma^2}\bigg)\mathbf{S}_{ija}\mathbf{S}_{ijb} \bigg\}
$$

For a given $i,j$, $\mathbf{S}_{ij}^T\mathbf{S}_{ij}$ gives us a matrix in $\mathbb{R}^{FH \times FH}$ whose $(a,b)$th element = $\mathbf{S}_{ija}\mathbf{S}_{ijb}$. The Hessian matrix is thus computed by

$$
\begin{equation}
\frac{\partial^2I_Q}{\partial\mathbf{w}^2} = \sum_i\sum_j\bigg\{ -\frac{\mathbf{M}_{ij}\mathbf{G}_{ij}}{2\sigma^2}\bigg(1 - \frac{\mathbf{Y}_{ij}^2}{2\sigma^2}\bigg)\mathbf{S}_{ij}^T\mathbf{S}_{ij} \bigg\}
\end{equation}
$$

### B. Gradient-based optimisation of quadratic mutual information

Equipped with the first and second derivatives, we can begin to consider applying gradient-based optimisation. The most commonly used optimisation algorithms are the Newton and quasi-Newton methods. Newton's method updates the guess at each iteration according to

$$
\begin{equation}
x_{k+1} = x_k - \alpha\bigg[\nabla^2 f(x_k)\bigg]^{-1}\nabla f(x_k)
\end{equation}
$$

quasi-Newton methods perform the same basic update, but substitute an approximation for $\nabla^2 f(x_k)$ where it is not avaialble. Newton and quasi-Newton methods yield fast and robust algorithms, but at the cost of poor scaling. If we're using a stimulus with frame size of 304 x 608, and a frame history of 5, then the Hessian matrix has shape $FH \times FH = 8.54 \times 10^{11}$ elements. This means that, at every iteration of the optimisation procedure, we would need to compute, store <i>and invert</i> a matrix with almost 1 trillion elements, which would effectively undermine the computational justification for using QMI over MI in the first place. Instead, we can turn to the simplest optimisation algorithm, the gradient method, which requires only knowledge of the jacobian and updates the guess at each iteration according to

$$
\begin{equation}
x_{k+1} = x_k - \alpha\nabla f(x_k)
\end{equation}
$$

Unfortunately, the cost of the gradient method's increased computational tractability is a significant decrease in speed of convergence - by considering only the first derivative, the gradient method assumes that the surface ot the objective function is planar, which will not be the case. This is not necessarily a major problem in the isolated case - but if we consider the need to batch-estimate the RFs of large populations of neurons, the poor speed of the algorithm compounds, and can become a serious limiting factor. Conjugate directions methods (CDM), first introduced by Hestenes and Steifel [[7](#7)], were developed to solve this tradeoff by offering, for quadratic functions, improved convergence speed without requiring computation of the Hessian. The update rule for the conjugate directions method is given by

$$
\begin{equation}
x_{k+1} = x_k + \alpha_kd_k
\end{equation}
$$

with

$$
\begin{equation}
\alpha = -\frac{d_k^T(\nabla^2f(x_k)x_k + \nabla f(x_k))}{d_k^T\nabla^2f(x_k)d_k}
\end{equation}
$$

and

$$
\begin{equation}
d_k = \begin{cases}
-\nabla f(x_0), & k = 0 \\
d_k = -\nabla f(x_k) + \beta_{k-1}d_{k-1}, & k \neq 0
\end{cases}
\end{equation}
$$

$$
\begin{equation}
\beta_k = \frac{\nabla f(x_{k+1})^T\nabla^2f(x_k)d_k}{d_k^T\nabla^2f(x_k)d_k}
\end{equation}
$$

At first glance, this doesn't seem to offer any improvement, since the Hessian $\nabla^2 f$ still appears multiple times in the update equations. The key insight here is that what we actually need is the ability fo use the Hessian to compute its product with a given vector $v$. To understand the importance of this distinction, we consider a single component of the product $(\nabla^2 f)v$. The $i$th row of $\nabla^2 f$ is a vector of partial derivatives of the form $\frac{\partial^2}{\partial x_i\partial x_j}f$, and so the $i$th row of $(\nabla^2 f)v$ is given by

$$
\begin{equation}
((\nabla^2 f)v)_i = \sum_{j=1}^N \frac{\partial^2 f}{\partial x_i\partial x_j}(x) \cdot v_j = \nabla \frac{\partial f}{\partial x_i}(x) \cdot v
\end{equation}
$$

This is simply the directional derivative of $\frac{\partial f}{\partial x_i}$ in the direction $v$. The definition of the directional derivative is given as

$$
\begin{equation}
\nabla_vf = \lim_{\epsilon \to 0}\frac{f(x + \epsilon v) - f(x)}{\epsilon}
\end{equation}
$$

which can be approximated using finite differences as

$$
\begin{equation}
\nabla_vf \approx \frac{f(x + \epsilon v) - f(x)}{\epsilon}
\end{equation}
$$

for some small $\epsilon$. We can thus approximate $\nabla^2f(x)v$ as

$$
\begin{equation}
\nabla^2f(x)v \approx \frac{\nabla f(x + \epsilon v) - \nabla f(x)}{\epsilon}
\end{equation}
$$

allowing us to replace each computation of the Hessian with two computations of the Jacobian.

### C. Implementation

[JAX](https://github.com/google/jax) [[8](#8)] is a python library developed by Google that offers GPU/TPU-accelerated equivalents of core NumPy and SciPy functionality, as well as a powerful autograd system that provides fast automatic differentiation of any native JAX function. For our purposes, the two most immediately relevant functions are <code>grad</code>, which takes a function $f$ and returns another function approximating its derivative $\nabla f$, and `hvp`, which takes a function $f$, a vector $v$, and a point $x$, and returns an approximation of the Hessian vector product described above. Combining these, it is possible to construct a very simple implementation of the CDM iteration:

```python
def alpha(fun, x, d):
    num = dot(d, grad(fun)(x))
    denom = dot(d, hvp(fun, x, d))

    return -(num / denom)

def beta(fun, x, d, j):
    Hd = hvp(fun, x, d)

    num = dot(j, Hd)
    denom = dot(d, Hd)

    return num / denom

while not converged:
  a = alpha(qmi, x, d)
  x += (a * d)
  x = x / norm(x)

  j = grad(fun)(x)
  b = beta(qmi, x, d, j)

  d = -j + (b * d)
```

<figcaption>Listing 1: JAX implementation of CDM loop</figcaption>

Here, <code>qmi</code> is a function that takes the current estimate of the receptive field vector, and returns the QMI between the stimulus projected onto that rfv, and the neural response vector. If initialised with a random receptive field vector with unit norm, this simple piece of code is <i>guaranteed</i> to converge in a finite number of iterations to the receptive field vector that yields the global maximum of <code>qmi</code>. Unfortunately, in the case of high-dimensional stimuli such as natural movies, this finite number can still be higher than we would ideally like, with the RF estimate remaining considerably noisy over a large number of iterations. In order to overcome this limitation, we can modify the optimisation approach in some way to <i>actively encourage</i> the denoising of the estimate, rather than just relying on some sufficiently large number of iterations to do it for us. There are two ways we might go about this - either modifying the objective function, thus keeping a single optimisation loop, or performing a separate noise removal optimisation after the main loop has reached a reasonable noisy estimate.

#### Modified objective function

Instead of maximising solely the QMI, we can modify the objective function to incorporate some measure of the <q>smoothness</q> of the current RF estimate. An intuitive candidate would be to compute some quantity based on the spatial frequency spectrum of the current receptive field vector. Specifically, we can compute the average spatial frequency of the RFV as

$$
\begin{equation}
\bar{\lambda} = \sum_{i=0}^{\text{len}(\mathbf{w})}i \cdot F[i]
\end{equation}
$$

where $F$ is the Fourier magnitude spectrum $|\mathcal{F}\{\mathbf{w}\}|$ normalised to lie within $[0,1]$. The smoothness can then be given as $1/\bar{\lambda}$, implemented straightforwardly in JAX as

```python
def smoothness(x):
    spectrum = normalise(abs(fft(x)))
    freq = average(range(len(spectrum)), weights=spectrum)
    return 1 / freq
```

<figcaption>Listing 2: JAX implementation of smoothness modifier</figcaption>

There are numerous ways that this measure could be incorporated into the existing objective function. The selection tested in this work were

$$
\begin{equation}
O_1(w) = I_q(w) \cdot \text{smoothness}(w)
\end{equation}
$$

$$
\begin{equation}
O_2(w) = I_q(w) + \gamma \cdot \text{smoothness}(w), \ \gamma > 0
\end{equation}
$$

$$
\begin{equation}
O_3(w) = I_q(w) \cdot \big[1 + \gamma\cdot\text{smoothness}(w)\big], \ \gamma > 0
\end{equation}
$$

#### Separate optimisation loop

The other possible approach is to denoise the RF estimate in a separate second loop. Given the noisy output of the initial QMI maximisation, we can perform progressive convolutions with a small averaging filter until the QMI stops increasing, and begins instead to decrease:

```python
rf, qmi_values = maximise_qmi(...)

w = 5
kernel = ones((w, w)) / (w ** 2)

while True:
  smoothed = convolve2d(rf, kernel) ** (1 + w/100)
  q = qmi(smoothed)

  if q <= max(qmi_values)
    break

  qmi_values.append(q)
  rf = smoothed
```

### D. Testing against simple model neurons

In order to validate the results of this optimisation procedure, we need to try estimating a known <q>ground-truth</q> RF, using a synthetic neural response. To generate the response, we return to the Poisson model described briefly in section [1A](#a-receptive-fields-and-their-estimation). For the first ground-truth RF, we use a purely spatial, 1-frame circular Gaussian filter:

<figure>
  <img src="./figures/template-1.svg">
  <figcaption>Fig 1. The spatial receptive field of the first model neuron used for testing</figcaption>
</figure>

The natural movie stimulus was then projected onto this filter, and the filter output at each time bin was squared. This produced a vector representing the number of spikes produced for each time bin (stimulus frame) by our model neuron. The stimulus and synthetic response were then used to generate RF estimates using both the optimisation algorithm and a simple implementation of the STA method outlined in section [1A](#a-receptive-fields-and-their-estimation), which were compared to the ground-truth. The same test was repeated with a second ground-truth RF incorporating a temporal variation over 3 frames:

<figure>
  <img src="./figures/template-2.svg">
  <figcaption>Fig 2. The spatiotemporal receptive field of the second model neuron used for testing</figcaption>
</figure>

### E. Application to the Allen Visual Coding dataset

To provide an example of how PyQMI's RF estimator might be used in computational neuroscience research, we used it to analyse neural data from [the Allen Institute's Large-Scale Visual Coding dataset](https://portal.brain-map.org/explore/circuits/visual-coding-neuropixels) (details of the dataset can be found in Appendix A). After loading all the data from session 715093703 (the first session in the dataset), an RF estimate was computed for each unit in the session. The full range of probe vertical positions from the session's units was determined, and divided into bins of width 100. For each bin, the average RF estimate for all units corresponding to the bin was determined. A pairwise similarity matrix was then produced from these average RFs, with the similarity defined as

$$
\begin{equation}
\text{similarity}(\mathbf{w}_i, \mathbf{w}_j) = 1 - \text{RMSE}(\mathbf{w}_i, \mathbf{w}_j)
\end{equation}
$$

and plotted as a heatmap. A similar analysis was then performed using the unit position in the left-right plane, rather than the vertical probe position.

```python
units = session.units
positions = units["probe_vertical_position"].values

bin_width = 100
num_bins = ceil((max(positions) - min(positions)) / bin_width)

rf_dict = {}
for i in range(num_bins):
    uids = units[
        units["probe_vertical_position"]
        in range(i * bin_width, (i + 1) * bin_width)
    ].index.values
    rfs = [estimate_rf_from_uid(uid) for uid in uids]
    rf_dict[i] = mean(rfs)

heatmap = np.zeros((num_bins, num_bins))
for i in range(num_bins):
  for j in range(num_bins):
    heatmap[i,j] = similarity(rf_dict[i], rf_dict[j])
```

<figcaption>Listing 3: Code for producing RF similarity matrix</figcaption>

## 3. Results

### A. Optimisation approaches

#### Comparing algorithms

<figure>
  <img src="./figures/qmi_curves.svg">
  <figcaption>Fig 3. Change in QMI over 50 iterations of the gradient and conjugate directions algorithms</figcaption>
</figure>

Fig. 3 shows the difference in convergence speed between the gradient and conjugate directions algorithms over 50 iterations of QMI optimisation (for the natural movie stimulus and a synthetic response produced by the first model neuron). It is evident that, while both produce QMI values that increase monotonically to the same maximum, the conjugate directions algorithm achieves faster convergence.

#### Incorporating RF smoothness

<row>
  <img src="figures/noisy.png" width="300"/> 
  <img src="figures/modified_objective_1.png" width="300"/>
</row>
<row>
  <img src="figures/modified_objective_2.png" width="300"/> 
  <img src="figures/smooth.png" width="300"/>
</row>

Fig 4. The results produced by different strategies of smoothness incorporation. Top left: the noisy RFV estimate produced by unmodified QMI optimisation. Top right: the result of optimising $O_3(w) = I_q(w) \cdot [1 + 100 \cdot \text{smoothness}(w)]$. Bottom left: the result of optimising $O_2(w) = I_q(w) + \text{smoothness}(w)$. Bottom right: the result of performing separate smoothing following the main optimisation loop.

As described in section [2C](#c-implementation), two broad approaches to incorporating RF smoothness were attempted. First, three modified objective functions were explored. Attempts at maximising $O_1 = I_q \cdot \text{smoothness}(w)$ failed to converge whatsoever. $O_2$ and $O_3$ did both yield convergence, but the results bear zero resemblance to the ground truth RF. This is likely because the signal from the QMI portion of each objective function was drowned out as the smoothness increased more quickly. The second approach was to target the smoothness of the RF estimate in a separate second loop, and the result can be seen in the bottom right of Figure 4. Clearly, this second approach is superior, because it avoids the problems and instabilities that arise from trying to balance the QMI and smoothness in a single objective function.

### B. Recovery of model neuron RFs

<figure>
<img src="./figures/comparison-1.png" width="400px">
<br>
<img src="./figures/comparison-2.png" width="400px">
<br>
<img src="./figures/comparison-3.png" width="400px">
<figcaption>Fig 5. Comparison between STA and QMI estimates of a model RF. Top: ground-truth RF; middle: STA estimate; bottom: QMI estimate</figcaption>
</figure>

Fig. 5 shows a comparison between the performance of the STA and QMI methods on recovery of a simple spatial RF from a Poisson model neuron. The STA method performs very poorly, due to the natural movie stimulus being highly non-Gaussian, with significant spatiotemporal correlation. The QMI method is not perfect, producing a not-entirely-spherical estimate with slightly overestimated radius, but it has clearly done a good job of recovering the essential structure of the ground-truth RF.

<figure>
    <img src="./figures/spatiotemporal.png" width="400px">
    <figcaption>Fig 6. Recovery of 3-frame spatiotemporal RF from model neuron by QMI optimiser</figcaption>
</figure>

Fig. 6 illustrates that the QMI method is not limited to purely spatial RFs, but can also recover spatiotemporal RFs at the cost of only slightly increased noise.

### C. Analysis of Allen dataset

Having validated the estimator on model neurons, we applied it to mapping a selection of units from the Allen dataset. Fig. 7 shows, as an example, the estimated RFs of two dLGN neurons taken from the first session of the dataset. Although these estimates are difficult to assess on account of not having access to ground truth values, we note that they do exhibit the typical centre-surround structure that would be expected of a neuron from the dLGN.

<figure>
    <img src="./figures/dlgn.png" width="400px">
    <figcaption>Fig 7. Estimated RFs of two dLGN neurons taken from the Allen dataset</figcaption>
</figure>

In order to demonstrate the applicability of the QMI method to analyses over neural populations, rather than just single cells, we produced two simple heatmaps. The first of these, given in Fig. 8, illustrates how the similarity between the average RFs for neural populations recorded by two given probes changes as a function of the vertical separation between the probes. As we might expect, in general, the similarity decreases as the probes become further vertically separated from one another. Fig. 9 shows the presence of a similar albeit weaker trend for separation along the left-right axis.

<figure>
  <img src="./figures/vertical.svg">
  <figcaption>Fig 8. The vertical probe position of every unit in the session was determined. For each bin of height 100, the average receptive field of all neurons in that bin was computed. The similarity between every pair of average RFs was then computed as 1 - RMSE, and plotted as a heatmap.</figcaption>
</figure>

<figure>
  <img src="./figures/left-right.svg">
  <figcaption>Fig 9. The position in the left-right plane for every unit in the session was determined. For each bin of width 100, the average receptive field of all neurons in that bin was computed. The similarity between every pair of average RFs was then computed as 1 - RMSE, and plotted as a heatmap.</figcaption>
</figure>

## 4. Discussion

### The need for better RF estimation tooling

Despite no longer being fit for purpose, the methods of spike-triggered analysis remain in popular and widespread use for receptive field estimation due to both their simplicity and their ubiquity of implementation. Although prior work has been carried out on RF estimation through QMI optimisation, it has not yet translated into the existence of a high-quality open-source software framework that could serve as a superior drop-in replacement for existing spike-triggered analysis libraries in general neuroinformatics workflows. The availability of such a framework would enable researchers and practitioners within the field to achieve greater insight into neural function with less data, as well as opening up the feasibility of accurately interrogating responses to less "synthetic", more biologically meaningful stimuli.

### Future work

Having demonstrated an effective implementation of a QMI-based RF estimator in GPU-accelerated python, the primary next step will be to ensure that the power of the method is accessible through an interface that simple and intuitive to use. The existing code was written largely with Allen dataset in mind - additional work can be done to enable easy and frictionless loading of data from multiple different sources. It is also likely that the code can be made even more performant, by taking greater advantage of the various low-level features and transforms offered by the JAX library. Finally, further investigations should be made into the optimisation logic itself. When modifying the approach of pure QMI maximisation, we have considered only the RF's quality of "smoothness", but it is possible that there may be other measures or attributes expected of reasonable RF estimates that can help guide the optimiser through the search space.

### Alternative approaches

Although RF estimation methods based on mutual information optimisation represent a meaningful and significant improvement on traditional approaches, they are by no means without limitations of their own; nor are they the only avenue currently being explored. In [[11](#11)], the authors discuss a possible limitation with optimization-based approaches to RF mapping, namely that they place the computational burden of searching for optimal parameters \textit{after} the collection of neural data has been completed, prohibiting their use in real-time closed-loop experiments. The authors aim to overcome this limitation by developing an alternate deep-learning-based approach that formulates the RF estimation procedure as a K-shot regression problem. It would be interesting to see if any elements of their approach could be used to adapt the QMI-based estimator to learn more efficiently from limited data, and potentially open up the possibility of real-time closed-loop usage.

## References

<span id='1'>[1] Sherrington, C S (1906). The integrative action of the nervous system. C Scribner and Sons, New York.</span>
<br>
<span id='2'>[2] Meister M, Pine J, Baylor DA. Multi-neuronal signals from the retina: acquisition and analysis. J Neurosci Methods. 1994 Jan;51(1):95-106.</span>
<br>
<span id='3'>[3] Reid RC, Alonso JM. Specificity of monosynaptic connections from thalamus to visual cortex. Nature. 1995 Nov 16;378(6554):281-4.</span>
<br>
<span id='4'>[4] McLean J, Palmer LA. Contribution of linear spatiotemporal receptive field structure to velocity selectivity of simple cells in area 17 of cat. Vision Res. 1989;29(6):675-9.</span>
<br>
<span id='5'>[5] Sharpee T, Rust NC, Bialek W. Analyzing neural responses to natural signals: maximally informative dimensions. Neural Comput. 2004 Feb;16(2):223-50</span>
<br>
<span id='6'>[6] Torkkola, K.. “Feature Extraction by Non-Parametric Mutual Information Maximization.” J. Mach. Learn. Res. 3 (2003): 1415-1438.</span>
<br>
<span id='7'>[7] Hestenes, M. and E. Stiefel. “Methods of conjugate gradients for solving linear systems.” Journal of research of the National Bureau of Standards 49 (1952): 409-435.</span>
<br>
<span id='8'>[8] Frostig, Roy, M. Johnson and Chris Leary. “Compiling machine learning programs via high-level tracing.” (2018).</span>
<br>
<span id='9'>[9] Katz, Matthew L., T. Viney and K. Nikolic. “Receptive Field Vectors of Genetically-Identified Retinal Ganglion Cells Reveal Cell-Type-Dependent Visual Functions.” PLoS ONE 11 (2016).</span>
<br>
<span id='10'>[10] Mu, Zhiguang, K. Nikolic and S. Schultz. “Quadratic Mutual Information estimation of mouse dLGN receptive fields reveals asymmetry between ON and OFF visual pathways.” bioRxiv (2020).</span>
<br>
<span id='11'>[11] Cotton, R. J., Fabian H. Sinz and A. Tolias. “Factorized Neural Processes for Neural Processes: $K$-Shot Prediction of Neural Responses.” ArXiv abs/2010.11810 (2020).
</span>
<br>
