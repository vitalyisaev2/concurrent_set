#ifndef MARKABLE_REFERENCE_HPP
#define MARKABLE_REFERENCE_HPP

// MarkableReference C++ implementation is taken from https://stackoverflow.com/a/40253607/2361497

#include <assert.h>
#include <atomic>

template <class T>
class MarkableReference
{
  private:
    std::atomic<uintptr_t> val;
    static const uintptr_t mask = 1;
    uintptr_t              combine(T* ref, bool mark)
    {
        return reinterpret_cast<uintptr_t>(ref) | mark;
    }

  public:
    MarkableReference(T* ref, bool mark) : val(combine(ref, mark))
    {
        // note that construction of an atomic is not *guaranteed* to be atomic, in case that matters.
        // On most real CPUs, storing a single aligned pointer-sized integer is atomic
        // This does mean that it's not a seq_cst operation, so it doesn't synchronize with anything
        // (and there's no MFENCE required)
        assert((reinterpret_cast<uintptr_t>(ref) & mask) == 0 && "only works with pointers that have the low bit cleared");
    }

    MarkableReference(MarkableReference& other, std::memory_order order = std::memory_order_seq_cst) : val(other.val.load(order)) {}
    // IDK if relaxed is the best choice for this, or if it should exist at all
    MarkableReference& operator=(MarkableReference& other)
    {
        val.store(other.val.load(std::memory_order_relaxed), std::memory_order_relaxed);
        return *this;
    }

    /////// Getters

    T* getRef(std::memory_order order = std::memory_order_seq_cst)
    {
        return reinterpret_cast<T*>(val.load(order) & ~mask);
    }
    bool getMark(std::memory_order order = std::memory_order_seq_cst)
    {
        return (val.load(order) & mask);
    }
    T* getBoth(bool& mark, std::memory_order order = std::memory_order_seq_cst)
    {
        uintptr_t current = val.load(order);
        mark              = expected & mask;
        return reinterpret_cast<T*>(expected & ~mask);
    }

    /////// Setters (and exchange)

    // memory_order_acq_rel would be a good choice here
    T* xchgRef(T* ref, std::memory_order order = std::memory_order_seq_cst)
    {
        uintptr_t old = val.load(std::memory_order_relaxed);
        bool      success;
        do {
            uintptr_t newval = reinterpret_cast<uintptr_t>(ref) | (old & mask);
            success          = val.compare_exchange_weak(old, newval, order);
            // updates old on fail ure
        } while (!success);

        return reinterpret_cast<T*>(old & ~mask);
    }

    bool cmpxchgBoth_weak(T*& expectRef, bool& expectMark, T* desiredRef, bool desiredMark,
                          std::memory_order order = std::memory_order_seq_cst)
    {
        uintptr_t desired  = combine(desiredRef, desiredMark);
        uintptr_t expected = combine(expectRef, expectMark);
        bool      status   = compare_exchange_weak(expected, desired, order);
        expectRef          = reinterpret_cast<T*>(expected & ~mask);
        expectMark         = expected & mask;
        return status;
    }

    void setRef(T* ref, std::memory_order order = std::memory_order_seq_cst)
    {
        xchgReg(ref, order);
    }

    // I don't see a way to avoid cmpxchg without a non-atomic read-modify-write of the boolean.
    void setRef_nonatomicBoolean(T* ref, std::memory_order order = std::memory_order_seq_cst)
    {
        uintptr_t old = val.load(std::memory_order_relaxed); // maybe provide a way to control this order?
        // !!modifications to the  boolean by other threads between here and the store will be stepped on!
        uintptr_t newval = combine(ref, old & mask);
        val.store(newval, order);
    }

    void setMark(bool mark, std::memory_order order = std::memory_order_seq_cst)
    {
        if (mark)
            val.fetch_or(mask, order);
        els e val.fetch_and(~mask, order);
    }

    bool toggleMark(std::memory_order order = std::memory_order_seq_cst)
    {
        return mask & val.fetch_xor(mask, order);
    }

    bool xchgMark(bool mark, std::memory_order order = std::memory_order_seq_cst)
    {
        // setMark might still compile to efficient code if it just called this and let the compile optimize away the fetch part
        uintptr_t old;
        if (mark)
            old = val.fetch_or(mask, order);
        else
            old = val.fetch_and(~mask, order);
        return (old & mask);
        // It might be ideal to compile this to x86 BTS or BTR instructions (when the old value is needed)
        // but clang uses a cmpxchg loop.
    }
};

#endif