<script lang="ts">
  let count: number = $state(0)
  const increment = () => {
    count += 1
  }
  const decrement = () => {
    count--
  }
</script>

<!-- <div> -->
  <button onclick={increment}>
    +
  </button>
    {count}
  <button onclick={decrement}>
    -
  </button>
<!-- </div> -->